package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type redemptionHandler struct {
	effectsProvider
}

func (h *redemptionHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	redemptionRequest := details.Request.Body.MustRedemptionRequest()
	redemptionRequestResponse := details.Result.TypeExt.MustCreateRedemptionResult().MustRedemptionResponse()

	source := h.effectsProvider.BalanceEffect(redemptionRequest.SourceBalanceId,
		&history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(redemptionRequest.Amount),
			},
		})

	destination := h.effectsProvider.BalanceEffect(redemptionRequestResponse.DestinationBalanceId,
		&history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: regources.Amount(redemptionRequest.Amount),
			},
		})

	return []history2.ParticipantEffect{source, destination}, nil
}

func (h *redemptionHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	redemptionRequest := details.Request.Body.MustRedemptionRequest()
	unlock := h.effectsProvider.BalanceEffect(redemptionRequest.SourceBalanceId,
		&history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(redemptionRequest.Amount),
			},
		})
	return []history2.ParticipantEffect{unlock}, nil
}
