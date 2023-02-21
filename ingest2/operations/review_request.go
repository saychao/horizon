package operations

import (
	"encoding/hex"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type reviewRequestOpHandler struct {
	effectsProvider
	allRequestHandlers map[xdr.ReviewableRequestType]reviewRequestHandler
}

type requestDetails struct {
	Op              xdr.ReviewRequestOp
	SourceAccountID xdr.AccountId
	Result          xdr.ExtendedResult
	Request         xdr.ReviewableRequestEntry
	Changes         []xdr.LedgerEntryChange
}

type reviewRequestHandler interface {
	//Fulfilled - returns participant of fully approved request
	Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error)
	//PermanentReject - returns participants of fully rejected request
	PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error)
}

// newReviewRequestOpHandler creates new handler for review request operations
// with specific handlers for different types
func newReviewRequestOpHandler(provider effectsProvider, defPayments defPaymentProvider) *reviewRequestOpHandler {
	// All reviewable requests must be explicitly handled here. If there is nothing to handle
	// use reviewableRequestHandlerStub
	stubProvider := reviewableRequestHandlerStub{
		effectsProvider: provider,
	}
	return &reviewRequestOpHandler{
		effectsProvider: provider,
		allRequestHandlers: map[xdr.ReviewableRequestType]reviewRequestHandler{
			xdr.ReviewableRequestTypeCreateIssuance: &issuanceHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeCreateWithdraw: &withdrawHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeCreateAmlAlert: &amlAlertHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeCreateAtomicSwapBid: &atomicSwapHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeCreateAsset:         &stubProvider,
			xdr.ReviewableRequestTypeUpdateAsset:         &stubProvider,
			xdr.ReviewableRequestTypeCreatePreIssuance:   &stubProvider,
			xdr.ReviewableRequestTypeCreateSale:          &stubProvider,
			xdr.ReviewableRequestTypeUpdateLimits:        &stubProvider,
			xdr.ReviewableRequestTypeChangeRole:          &stubProvider,
			xdr.ReviewableRequestTypeUpdateSaleDetails:   &stubProvider,
			xdr.ReviewableRequestTypeCreateInvoice:       &deprecatedReviewRequestHandler{},
			xdr.ReviewableRequestTypeManageContract:      &deprecatedReviewRequestHandler{},
			xdr.ReviewableRequestTypeCreateAtomicSwapAsk: &stubProvider,
			xdr.ReviewableRequestTypeCreatePoll:          &stubProvider,
			xdr.ReviewableRequestTypeKycRecovery:         &stubProvider,
			xdr.ReviewableRequestTypeManageOffer: &manageOfferRequestHandler{
				&manageOfferOpHandler{provider}},
			xdr.ReviewableRequestTypeCreatePayment: &paymentRequestHandler{
				&paymentOpHandler{provider},
			},
			xdr.ReviewableRequestTypePerformRedemption: &redemptionHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeDataCreation: &createDataHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeDataUpdate: &updateDataHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeDataRemove: &removeDataHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeCreateDeferredPayment: &createDeferredPaymentHandler{
				effectsProvider: provider,
			},
			xdr.ReviewableRequestTypeCloseDeferredPayment: &closeDeferredPaymentHandler{
				effectsProvider:    provider,
				defPaymentProvider: defPayments,
			},
		},
	}
}

// Details returns details about review request operation
func (h *reviewRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	reviewRequestOp := op.Body.MustReviewRequestOp()
	reviewRequestOpRes := opRes.MustReviewRequestResult().MustSuccess()

	addedTasks := uint32(reviewRequestOp.ReviewDetails.TasksToAdd)
	removedTasks := uint32(reviewRequestOp.ReviewDetails.TasksToRemove)
	externalDetails := internal.MarshalCustomDetails(xdr.Longstring(reviewRequestOp.ReviewDetails.ExternalDetails))

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeReviewRequest,
		ReviewRequest: &history2.ReviewRequestDetails{
			RequestID:       int64(reviewRequestOp.RequestId),
			RequestType:     reviewRequestOp.RequestDetails.RequestType,
			RequestHash:     hex.EncodeToString(reviewRequestOp.RequestHash[:]),
			Action:          reviewRequestOp.Action,
			Reason:          string(reviewRequestOp.Reason),
			IsFulfilled:     reviewRequestOpRes.Fulfilled,
			RequestDetails:  reviewRequestOp.RequestDetails,
			AddedTasks:      addedTasks,
			RemovedTasks:    removedTasks,
			ExternalDetails: externalDetails,
		},
	}

	return opDetails, nil
}

// ParticipantsEffects - returns source participant if request was not fulfilled
// finds specific handler otherwise
func (h *reviewRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	reviewRequestOp := opBody.MustReviewRequestOp()
	reviewRequestRes := opRes.MustReviewRequestResult().MustSuccess()
	request := h.getReviewableRequestByID(int64(reviewRequestOp.RequestId), ledgerChanges)

	if request == nil {
		return nil, errors.From(
			errors.New("expected request to be in ledger changes"), map[string]interface{}{
				"request_id": int64(reviewRequestOp.RequestId),
			})
	}

	specificHandler, ok := h.allRequestHandlers[request.Body.Type]
	if !ok {
		return nil, errors.From(errors.New("failed to find handler for reviewable request"), logan.F{
			"request_type": request.Body.Type,
		})
	}

	details := requestDetails{
		Op:              reviewRequestOp,
		SourceAccountID: sourceAccountID,
		Result:          reviewRequestRes,
		Request:         *request,
		Changes:         ledgerChanges,
	}

	switch reviewRequestOp.Action {
	case xdr.ReviewRequestOpActionApprove:
		if !reviewRequestRes.Fulfilled {
			return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
		}

		return specificHandler.Fulfilled(details)
	case xdr.ReviewRequestOpActionReject:
		return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
	case xdr.ReviewRequestOpActionPermanentReject:
		return specificHandler.PermanentReject(details)
	default:
		return nil, errors.From(errors.New("unexpected action on reivew request participants"), logan.F{
			"action": reviewRequestOp.Action,
		})

	}
}

// Tries to get latest version of reviewable request by ID (updated first, created or state otherwise)
func (h *reviewRequestOpHandler) getReviewableRequestByID(id int64, ledgerChanges []xdr.LedgerEntryChange,
) *xdr.ReviewableRequestEntry {
	var bestResult *xdr.ReviewableRequestEntry

	for _, change := range ledgerChanges {
		var ledgerEntryData xdr.LedgerEntryData

		switch change.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			ledgerEntryData = change.MustCreated().Data
		case xdr.LedgerEntryChangeTypeUpdated:
			ledgerEntryData = change.MustUpdated().Data
		case xdr.LedgerEntryChangeTypeState:
			ledgerEntryData = change.MustState().Data
		default:
			continue
		}

		if ledgerEntryData.Type != xdr.LedgerEntryTypeReviewableRequest {
			continue
		}

		request := ledgerEntryData.MustReviewableRequest()
		if int64(request.RequestId) != id {
			continue
		}

		bestResult = &request
		// we have found latest version of the request, so can return
		if change.Type == xdr.LedgerEntryChangeTypeUpdated {
			return bestResult
		}
	}

	return bestResult
}
