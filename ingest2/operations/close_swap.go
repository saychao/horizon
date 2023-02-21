package operations

import (
	"encoding/hex"

	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"

	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type closeSwapOpHandler struct {
	effectsProvider
	swapProvider
}

// Details returns details about manage balance operation
func (h *closeSwapOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	closeSwapOp := op.Body.MustCloseSwapOp()
	res := opRes.MustCloseSwapResult()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCloseSwap,
		CloseSwap: &history2.CloseSwapDetails{
			ID: int64(closeSwapOp.SwapId),
		},
	}

	if res.MustSuccess().Effect == xdr.CloseSwapEffectClosed {
		secret := hex.EncodeToString(closeSwapOp.Secret[:])
		details.CloseSwap.Secret = &secret
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *closeSwapOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return h.participantEffects(opBody.MustCloseSwapOp(),
		opRes.MustCloseSwapResult().MustSuccess(),
		sourceAccountID)
}

func (h *closeSwapOpHandler) participantEffects(op xdr.CloseSwapOp,
	res xdr.CloseSwapSuccess, sourceAccountID xdr.AccountId,
) ([]history2.ParticipantEffect, error) {

	switch res.Effect {
	case xdr.CloseSwapEffectClosed:
		return h.getClosedParticipants(op, res, sourceAccountID)
	case xdr.CloseSwapEffectCancelled:
		return h.getCancelledParticipants(op, res, sourceAccountID)
	default:
		return nil, errors.From(errors.New("unexpected close swap effect"), map[string]interface{}{
			"effect":  res.Effect.String(),
			"swap_id": uint64(op.SwapId),
		})
	}

}

func (h *closeSwapOpHandler) getClosedParticipants(op xdr.CloseSwapOp,
	res xdr.CloseSwapSuccess, sourceAccountID xdr.AccountId) ([]history2.ParticipantEffect, error) {
	swap := h.MustSwap(int64(op.SwapId))

	var sourceBalance xdr.BalanceId
	err := sourceBalance.SetString(swap.SourceBalance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse balance id from string", logan.F{
			"swap_id": int64(op.SwapId),
		})
	}

	var destBalance xdr.BalanceId
	err = destBalance.SetString(swap.DestinationBalance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse balance id from string", logan.F{
			"swap_id": int64(op.SwapId),
		})
	}

	source := h.BalanceEffect(sourceBalance, &history2.Effect{
		Type: history2.EffectTypeChargedFromLocked,
		ChargedFromLocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(swap.Amount),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(swap.SourcePercentFee),
				Fixed:             regources.Amount(swap.SourceFixedFee),
			},
		},
	})

	dest := h.BalanceEffect(destBalance, &history2.Effect{
		Type: history2.EffectTypeFunded,
		Funded: &history2.BalanceChangeEffect{
			Amount: regources.Amount(swap.Amount),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(swap.DestinationPercentFee),
				Fixed:             regources.Amount(swap.DestinationFixedFee),
			},
		},
	})
	return []history2.ParticipantEffect{source, dest}, nil
}

func (h *closeSwapOpHandler) getCancelledParticipants(op xdr.CloseSwapOp,
	res xdr.CloseSwapSuccess, sourceAccountID xdr.AccountId) ([]history2.ParticipantEffect, error) {
	swap := h.MustSwap(int64(op.SwapId))

	var sourceBalance xdr.BalanceId
	err := sourceBalance.SetString(swap.SourceBalance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse balance id from string", logan.F{
			"swap_id": int64(op.SwapId),
		})
	}

	source := h.BalanceEffect(sourceBalance, &history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(swap.Amount),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(swap.SourcePercentFee + swap.DestinationPercentFee),
				Fixed:             regources.Amount(swap.SourceFixedFee + swap.DestinationFixedFee),
			},
		},
	})
	return []history2.ParticipantEffect{source}, nil
}
