package handlers

import (
	"net/http"
	"strconv"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources/generated"
)

// GetTransaction - processes request to get the transaction by id or hash
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getTransactionHandler{
		TransactionsQ:  history2.NewTransactionsQ(historyRepo),
		LedgerChangesQ: history2.NewLedgerChangesQ(historyRepo),
		LedgerQ:        *history2.NewLedgerQ(historyRepo),
		Log:            ctx.Log(r),
	}

	request, err := requests.NewGetTransaction(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	result, err := handler.GetTransaction(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get transactions list")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getTransactionHandler struct {
	TransactionsQ  history2.TransactionsQ
	LedgerChangesQ history2.LedgerChangesQ
	LedgerQ        history2.LedgerQ
	Log            *logan.Entry
}

func (h *getTransactionHandler) tryGetByID(idOrHash string) (*history2.Transaction, error) {
	// we might have received request to get by hash, so just return nil, nil
	id, err := strconv.ParseUint(idOrHash, 10, 64)
	if err != nil {
		return nil, nil
	}

	result, err := h.TransactionsQ.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction by id")
	}

	return result, nil
}

func (h *getTransactionHandler) tryGetTX(idOrHash string) (*history2.Transaction, error) {
	tx, err := h.tryGetByID(idOrHash)
	if err != nil {
		return nil, err
	}

	if tx != nil {
		return tx, nil
	}

	tx, err = h.TransactionsQ.FilterByHash(idOrHash).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx by hash")
	}

	return tx, nil
}

// GetTransactions returns the list of transactions with related resources
func (h *getTransactionHandler) GetTransaction(request *requests.GetTransaction) (*regources.TransactionResponse, error) {
	historyTx, err := h.tryGetTX(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx")
	}

	if historyTx == nil {
		return nil, nil
	}

	var result regources.TransactionResponse
	result.Data, err = getPopulatedTx(*historyTx, h.LedgerChangesQ, request, &result.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populated tx")
	}

	return &result, nil
}
