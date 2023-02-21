package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	history "github.com/SafeRE-IT/horizon/db2/history2"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetTransactions - processes request to get the list of transactions (with ledger changes)
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getTransactionsHandler{
		LedgerQ:        *history2.NewLedgerQ(historyRepo),
		LedgerChangesQ: history2.NewLedgerChangesQ(historyRepo),
		TransactionsQ:  history2.NewTransactionsQ(historyRepo),
		Log:            ctx.Log(r),
	}

	request, err := requests.NewGetTransactions(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	result, err := handler.GetTransactions(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get transactions list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getTransactionsHandler struct {
	TransactionsQ  history2.TransactionsQ
	LedgerChangesQ history2.LedgerChangesQ
	LedgerQ        history2.LedgerQ
	Log            *logan.Entry
}

func (h *getTransactionsHandler) getLatestLedger() (*history2.Ledger, error) {
	sequence, err := h.LedgerQ.GetLatestLedgerSeq()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load latest ledger sequence")
	}
	h.Log.WithField("latest_sequence", sequence).Debug("Got latest ledger seq")

	ledger, err := h.LedgerQ.GetBySequence(sequence)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger by sequence")
	}
	if ledger == nil {
		h.Log.WithField("latest_sequence", sequence).Error("Failed to get latest ledger")
	}

	return ledger, nil
}

// GetTransactions returns the list of transactions with related resources
func (h *getTransactionsHandler) GetTransactions(request *requests.GetTransactions) (*regources.TransactionListResponse, error) {
	q := h.TransactionsQ.Page(request.PageParams)

	if request.Filters.ChangeTypes != nil {
		q = q.FilterByEffectTypes(request.Filters.ChangeTypes...)
	}

	if request.Filters.EntryTypes != nil {
		q = q.FilterByLedgerEntryTypes(request.Filters.EntryTypes...)
	}

	if request.Filters.AfterTimestamp != nil {
		gottime := time.Unix(*request.Filters.AfterTimestamp, 0)
		fmt.Println(gottime)
		q = q.FilterLedgerCloseTimeAfter(time.Unix(*request.Filters.AfterTimestamp, 0).UTC())
	}

	if request.Filters.BeforeTimestamp != nil {
		gottime := time.Unix(*request.Filters.BeforeTimestamp, 0)
		fmt.Println(gottime)
		q = q.FilterLedgerCloseTimeBefore(time.Unix(*request.Filters.BeforeTimestamp, 0).UTC())
	}

	historyTransactions, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load transactions")
	}

	result := regources.TransactionListResponse{
		Data: make([]regources.Transaction, 0, len(historyTransactions)),
	}

	for _, historyTransaction := range historyTransactions {
		var transaction regources.Transaction
		transaction, err = getPopulatedTx(historyTransaction, h.LedgerChangesQ, request, &result.Included)
		if err != nil {
			return nil, errors.Wrap(err, "failed to populate tx")
		}

		result.Data = append(result.Data, transaction)
	}

	if len(result.Data) > 0 {
		result.Links = request.GetCursorLinks(request.PageParams, result.Data[len(result.Data)-1].ID)
	} else {
		result.Links = request.GetCursorLinks(request.PageParams, "")
	}

	latestLedger, err := h.getLatestLedger()
	// TODO: possible race condition may occur if new ledger will be closed between querying transactions and ledgers
	// need to find a solution and fix somehow
	if err != nil {
		return nil, errors.Wrap(err, "failed to load latest ledger")
	}
	if latestLedger == nil {
		return nil, errors.New("Latest ledger is nil")
	}

	if result.Meta, err = json.Marshal(regources.TransactionResponseMeta{
		LatestLedgerCloseTime: latestLedger.ClosedAt,
		LatestLedgerSequence:  latestLedger.Sequence,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to marshal meta")
	}

	return &result, nil
}

type baseRequest interface {
	ShouldInclude(string) bool
}

func getPopulatedTx(tx history.Transaction, ledgerChangesQ history.LedgerChangesQ, request baseRequest, include *regources.Included) (regources.Transaction, error) {
	historyChanges, err := ledgerChangesQ.FilterByTransactionID(tx.ID).OrderByNumber().Select()
	if err != nil {
		return regources.Transaction{}, errors.Wrap(err, "failed to load ledger changes")
	}

	transaction := resources.NewTransaction(tx)
	transaction.Relationships.Source = resources.NewAccountKey(tx.Account).AsRelation()
	transaction.Relationships.LedgerEntryChanges = &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(historyChanges)),
	}

	operations := make(map[int64]regources.Key)

	for _, historyChange := range historyChanges {
		var change *regources.LedgerEntryChange
		change, err = resources.NewLedgerEntryChange(historyChange)
		if err != nil {
			return regources.Transaction{}, errors.Wrap(err, "failed to parse ledger entry change")
		}
		transaction.Relationships.LedgerEntryChanges.Data = append(
			transaction.Relationships.LedgerEntryChanges.Data, change.Key,
		)

		if request.ShouldInclude(requests.IncludeTypeTransactionListLedgerEntryChanges) {
			include.Add(change)
		}
		operations[historyChange.OperationID] = resources.NewOperationKey(historyChange.OperationID)
	}

	transaction.Relationships.Operations = &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(operations)),
	}
	for _, operation := range operations {
		transaction.Relationships.Operations.Data = append(
			transaction.Relationships.Operations.Data, operation,
		)
	}

	return transaction, nil
}
