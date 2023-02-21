package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type payoutHandler struct {
	effectsProvider
}

// Details returns details about payout operation
func (h *payoutHandler) Details(op rawOperation, res xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	payoutOp := op.Body.MustPayoutOp()
	payoutRes := res.MustPayoutResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypePayout,
		Payout: &history2.PayoutDetails{
			SourceAccountAddress: op.Source.Address(),
			SourceBalanceAddress: payoutOp.SourceBalanceId.AsString(),
			Asset:                string(payoutOp.Asset),
			MaxPayoutAmount:      regources.Amount(payoutOp.MaxPayoutAmount),
			MinAssetHolderAmount: regources.Amount(payoutOp.MinAssetHolderAmount),
			MinPayoutAmount:      regources.Amount(payoutOp.MinPayoutAmount),
			ExpectedFee:          internal.FeeFromXdr(payoutOp.Fee),
			ActualFee:            internal.FeeFromXdr(payoutRes.ActualFee),
			ActualPayoutAmount:   regources.Amount(payoutRes.ActualPayoutAmount),
		},
	}, nil
}

// ParticipantsEffects returns `charged` and `funded` effects
func (h *payoutHandler) ParticipantsEffects(opBody xdr.OperationBody,
	res xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	payoutOp := opBody.MustPayoutOp()
	payoutRes := res.MustPayoutResult().MustSuccess()

	source := h.BalanceEffect(payoutOp.SourceBalanceId, &history2.Effect{
		Type: history2.EffectTypeCharged,
		Charged: &history2.BalanceChangeEffect{
			Amount: regources.Amount(payoutRes.ActualPayoutAmount),
			Fee:    internal.FeeFromXdr(payoutRes.ActualFee),
		},
	})

	responses := payoutRes.PayoutResponses
	participants := make([]history2.ParticipantEffect, 0, len(responses))
	participants = append(participants, source)
	for _, response := range responses {
		effect := history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: regources.Amount(response.ReceivedAmount),
			},
		}

		participants = append(participants, h.BalanceEffect(response.ReceiverBalanceId, &effect))
	}

	return participants, nil
}
