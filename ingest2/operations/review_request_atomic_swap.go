package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/regources/generated"
)

type atomicSwapHandler struct {
	effectsProvider
}

// PermanentReject - returns participants of fully rejected request
func (h *atomicSwapHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return h.tryHandleUnlockedEffect(details, []history2.ParticipantEffect{h.Participant(details.SourceAccountID)})
}

// Fulfilled - returns slice of effects for participants of the operation
func (h *atomicSwapHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	atomicSwapExtendedResult := details.Result.TypeExt.MustAtomicSwapBidExtended()

	participants := []history2.ParticipantEffect{
		h.BalanceEffect(atomicSwapExtendedResult.AskOwnerBaseBalanceId, &history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(atomicSwapExtendedResult.BaseAmount),
			},
		}),
	}

	participants = append(participants,
		h.BalanceEffect(atomicSwapExtendedResult.BidOwnerBaseBalanceId, &history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: regources.Amount(atomicSwapExtendedResult.BaseAmount),
			},
		}))

	return h.tryHandleUnlockedEffect(details, participants)
}

func (h *atomicSwapHandler) tryHandleUnlockedEffect(details requestDetails,
	participants []history2.ParticipantEffect,
) ([]history2.ParticipantEffect, error) {
	atomicSwapExtendedResult := details.Result.TypeExt.MustAtomicSwapBidExtended()

	if atomicSwapExtendedResult.UnlockedAmount == 0 {
		return participants, nil
	}

	participants = append(participants,
		h.BalanceEffect(atomicSwapExtendedResult.AskOwnerBaseBalanceId, &history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(atomicSwapExtendedResult.UnlockedAmount),
			},
		}))

	return participants, nil
}
