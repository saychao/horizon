package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createManageOfferRequestOpHandler struct {
	offerHandler *manageOfferOpHandler
}

// Details returns details about create limits request operation
func (h *createManageOfferRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createManageOfferRequestOp := op.Body.MustCreateManageOfferRequestOp()
	createManageOfferRequestRes := opRes.MustCreateManageOfferRequestResult().MustSuccess()
	request := createManageOfferRequestOp.Request
	manageOfferOp := createManageOfferRequestOp.Request.Op

	var allTasks *uint32
	if createManageOfferRequestOp.AllTasks != nil {
		tasks := uint32(*createManageOfferRequestOp.AllTasks)
		allTasks = &tasks
	}

	creatorDetails := regources.Details("{}")
	switch request.Ext.V {
	case xdr.LedgerVersionMovementRequestsDetails:
		creatorDetails = internal.MarshalCustomDetails(request.Ext.MustCreatorDetails())
	case xdr.LedgerVersionEmptyVersion:
	default:
		panic(errors.From(errors.New("unexpected version of payment request"), logan.F{
			"ledger_version": request.Ext.V,
		}))
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateManageOfferRequest,
		CreateManageOfferRequest: &history2.CreateManageOfferRequestDetails{
			ManageOfferDetails: history2.ManageOfferRequest{
				CreatorDetails: creatorDetails,
				OfferID:        int64(manageOfferOp.OfferId),
				OrderBookID:    int64(manageOfferOp.OrderBookId),
				Amount:         regources.Amount(manageOfferOp.Amount),
				Price:          regources.Amount(manageOfferOp.Price),
				IsBuy:          manageOfferOp.IsBuy,
				Fee: regources.Fee{
					CalculatedPercent: regources.Amount(manageOfferOp.Fee),
				},
			},
			AllTasks: allTasks,
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(createManageOfferRequestRes.RequestId),
				IsFulfilled: createManageOfferRequestRes.Fulfilled,
			},
		},
	}, nil
}

func (h *createManageOfferRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, changes []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustCreateManageOfferRequestOp().Request.Op
	result := opRes.MustCreateManageOfferRequestResult().MustSuccess()

	if !result.Fulfilled {
		return []history2.ParticipantEffect{h.offerHandler.Participant(sourceAccountID)}, nil
	}

	if manageOfferOp.Amount != 0 {
		source := h.offerHandler.Participant(sourceAccountID)
		if result.ManageOfferResult == nil {
			return nil, errors.New("unexpected nil manage offer result")
		}
		return h.offerHandler.getNewOfferEffect(manageOfferOp, result.ManageOfferResult.MustSuccess(), source, changes), nil
	}

	deletedOfferEffects := h.offerHandler.getDeletedOffersEffect(changes)
	if len(deletedOfferEffects) != 1 {
		return nil, errors.From(errors.New("Unexpected number of deleted offer for delete offer though request"), logan.F{
			"expected": 1,
			"actual":   len(deletedOfferEffects),
		})
	}

	return deletedOfferEffects, nil
}
