package horizon

import (
	"github.com/saychao/horizon/db2"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"golang.org/x/net/context"
)

// LedgerOperationsIndexAction is an Action-based struct designed to co
type LedgerOperationsIndexAction struct {
	Action
	Types        []xdr.OperationType
	PagingParams db2.PageQuery
	Records      resource.Data
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *LedgerOperationsIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *LedgerOperationsIndexAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
}

func collectLedgers(historyQ history.QInterface, pagingParams db2.PageQuery, ctx context.Context) (result []resource.DataLedger, err error) {
	var ledgers []history.Ledger
	err = historyQ.Ledgers().Page(pagingParams).Select(&ledgers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select ledgers by paging params")
	}
	for _, ledger := range ledgers {
		dataLedger := resource.DataLedger{
			ClosedAt:     ledger.ClosedAt,
			Sequence:     ledger.Sequence,
			LedgerHash:   ledger.LedgerHash,
			Transactions: nil,
		}
		dataLedger.Transactions, err = collectLedgerTransactions(historyQ, ledger, ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to collect ledger transactions")
		}
		result = append(result, dataLedger)
	}
	return result, nil
}

func collectLedgerTransactions(historyQ history.QInterface, ledger history.Ledger, ctx context.Context) (result []resource.DataLedgerTransaction, err error) {
	var transactions []history.Transaction
	err = historyQ.Transactions().ForLedger(int32(ledger.Sequence)).Select(&transactions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select transactions for ledger")
	}
	for _, transaction := range transactions {
		dataTransaction := resource.DataLedgerTransaction{
			ID:         transaction.ID,
			Operations: nil,
		}
		dataTransaction.Operations, err = collectTransactionOperations(historyQ, transaction, ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to collect operations")
		}
		result = append(result, dataTransaction)
	}
	return result, nil
}

func collectTransactionOperations(historyQ history.QInterface, transaction history.Transaction, ctx context.Context) (result []hal.Pageable, err error) {
	var operations []history.Operation
	err = historyQ.Operations().ForTx(transaction.TransactionHash).Select(&operations)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select operations for tx")
	}
	for _, operation := range operations {
		dataOperation, err := resource.NewPublicOperation(ctx, operation, nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to populate public operation")
		}
		result = append(result, dataOperation)
	}
	return result, nil
}

func (action *LedgerOperationsIndexAction) loadRecords() {
	var data resource.Data
	var err error
	data.Ledgers, err = collectLedgers(action.HistoryQ(), action.PagingParams, action.Ctx)
	if err != nil {
		action.Log.WithError(err).Error("failed to collect records")
		action.Err = &problem.ServerError
	}
	action.Records = data
}

func (action *LedgerOperationsIndexAction) loadPage() {
	for _, record := range action.Records.Ledgers {
		action.Page.Add(record)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
