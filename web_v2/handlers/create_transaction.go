package handlers

import (
	"context"
	"net/http"

	"github.com/google/jsonapi"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/saychao/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/saychao/horizon/db2/history2"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/txsub/v2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

// CreateTransaction submits transaction to the core
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateTransactionRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)

	handler := createTransactionHandler{
		Core:      core2.NewTransactionQ(coreRepo),
		History:   history2.NewTransactionsQ(historyRepo),
		Log:       ctx.Log(r),
		Submitter: ctx.Submitter(r),
	}

	handler.Log = handler.Log.WithFields(logan.F{
		"tx_hash": request.Env.ContentHash,
	})

	result, err := handler.createTx(r.Context(), request)
	if errObj, ok := err.(*jsonapi.ErrorObject); ok {
		ctx.Log(r).WithError(err).WithFields(logan.F{
			"request": request,
		}).Error("failed to create transaction ")
		ape.RenderErr(w, errObj)
		return
	}
	if err != nil {
		ctx.Log(r).WithError(err).WithFields(logan.F{
			"request": request,
		}).Error("failed to create transaction ")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ctx.Log(r).WithFields(logan.F{
			"request": request,
		}).Error("got empty result")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type createTransactionHandler struct {
	History   history2.TransactionsQ
	Core      core2.TransactionQ
	Submitter *txsub.System
	Log       *logan.Entry
}

func (h *createTransactionHandler) createTx(context context.Context, request *requests.CreateTransaction) (*regources.TransactionResponse, error) {
	res, err := h.Submitter.Submit(context, *request.Env, request.WaitForResult, request.WaitForIngest)
	if txsub.IsInternalError(errors.Cause(err)) {
		return nil, resources.NewTxFailure(*request.Env, errors.Cause(err).(txsub.Error))
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to handle create transaction request")
	}
	if res == nil {
		return nil, errors.New("failed to submit transaction")
	}

	if !request.WaitForResult {
		return h.prepareWithoutResultXDR(res)
	}

	if request.WaitForIngest {
		return h.prepareFromHistory(uint64(res.TransactionID))
	}

	return h.prepareFromResult(res)
}

func (h *createTransactionHandler) prepareFromHistory(ID uint64) (*regources.TransactionResponse, error) {
	tx, err := h.History.GetByID(ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx by id", logan.F{
			"tx_id": ID,
		})
	}
	if tx == nil {
		return nil, nil
	}

	response := &regources.TransactionResponse{}
	response.Data = resources.NewTransaction(*tx)
	return response, nil
}

func (h *createTransactionHandler) prepareWithoutResultXDR(result *txsub.Result) (*regources.TransactionResponse, error) {
	return &regources.TransactionResponse{
		Data: regources.Transaction{
			Key: resources.NewTxKeyFromHash(result.Hash),
			Attributes: regources.TransactionAttributes{
				EnvelopeXdr: result.EnvelopeXDR,
				Hash:        result.Hash,
			},
		},
	}, nil
}

func (h *createTransactionHandler) prepareFromResult(result *txsub.Result) (*regources.TransactionResponse, error) {
	response := &regources.TransactionResponse{}
	data := regources.Transaction{
		Key: resources.NewTxKey(result.TransactionID),
		Attributes: regources.TransactionAttributes{
			EnvelopeXdr:    result.EnvelopeXDR,
			Hash:           result.Hash,
			LedgerSequence: result.LedgerSequence,
			ResultMetaXdr:  result.ResultMetaXDR,
			ResultXdr:      result.ResultXDR,
		},
	}
	response.Data = data
	return response, nil
}
