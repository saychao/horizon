package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
)

type paymentRequestHandler struct {
	paymentHandler *paymentOpHandler
}

// PermanentReject - returns participants of fully rejected request
func (h *paymentRequestHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.paymentHandler.Participant(details.SourceAccountID)}, nil
}

// Fulfilled - returns slice of effects for participants of the operation
func (h *paymentRequestHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return h.paymentHandler.participantEffects(details.Request.Body.MustCreatePaymentRequest().PaymentOp,
		details.Result.TypeExt.MustPaymentResult().MustPaymentResponse(),
		details.Request.Requestor)
}
