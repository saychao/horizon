package operations

import (
	"github.com/saychao/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type createDeferredPaymentHandler struct {
	effectsProvider
}

func (h *createDeferredPaymentHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.effectsProvider.Participant(details.SourceAccountID)}, nil
}

func (h *createDeferredPaymentHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	createDeferredPaymentRequest := details.Request.Body.MustCreateDeferredPaymentRequest()

	unlocked := h.effectsProvider.BalanceEffect(createDeferredPaymentRequest.SourceBalance,
		&history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(createDeferredPaymentRequest.Amount),
			},
		})

	return []history2.ParticipantEffect{unlocked}, nil
}
