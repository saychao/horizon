package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type manageOfferRequestHandler struct {
	offerHandler *manageOfferOpHandler
}

// PermanentReject - returns participants of fully rejected request
func (h *manageOfferRequestHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.offerHandler.Participant(details.SourceAccountID)}, nil
}

// Fulfilled - returns slice of effects for participants of the operation
func (h *manageOfferRequestHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	manageOfferOp := details.Request.Body.MustManageOfferRequest().Op
	result := details.Result.TypeExt.MustManageOfferResult()

	if manageOfferOp.Amount != 0 {
		source := h.offerHandler.Participant(details.Request.Requestor)
		return h.offerHandler.getNewOfferEffect(manageOfferOp, result.MustSuccess(), source, details.Changes), nil
	}

	deletedOfferEffects := h.offerHandler.getDeletedOffersEffect(details.Changes)
	if len(deletedOfferEffects) != 1 {
		return nil, errors.From(errors.New("Unexpected number of deleted offer for delete offer though request"), logan.F{
			"expected": 1,
			"actual":   len(deletedOfferEffects),
		})
	}

	return deletedOfferEffects, nil
}
