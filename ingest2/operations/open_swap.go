package operations

import (
	"encoding/hex"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources/generated"
)

type openSwapOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *openSwapOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	openSwapOp := op.Body.MustOpenSwapOp()
	openSwapRes := opRes.MustOpenSwapResult().MustSuccess()

	secretHash := hex.EncodeToString(openSwapOp.SecretHash[:])

	details := history2.OperationDetails{
		Type: xdr.OperationTypeOpenSwap,
		OpenSwap: &history2.OpenSwapDetails{
			AccountFrom:             op.Source.Address(),
			AccountTo:               openSwapRes.Destination.Address(),
			BalanceFrom:             openSwapOp.SourceBalance.AsString(),
			BalanceTo:               openSwapRes.DestinationBalance.AsString(),
			Amount:                  regources.Amount(openSwapOp.Amount),
			Asset:                   string(openSwapRes.Asset),
			SourceFee:               internal.FeeFromXdr(openSwapRes.ActualSourceFee),
			DestinationFee:          internal.FeeFromXdr(openSwapRes.ActualDestinationFee),
			SourcePayForDestination: openSwapOp.FeeData.SourcePaysForDest,
			SecretHash:              secretHash,
			LockTime:                internal.TimeFromXdr(xdr.Uint64(openSwapOp.LockTime)),
			Details:                 nil,
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *openSwapOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return h.participantEffects(opBody.MustOpenSwapOp(),
		opRes.MustOpenSwapResult().MustSuccess(),
		sourceAccountID)
}

func (h *openSwapOpHandler) participantEffects(op xdr.OpenSwapOp,
	res xdr.OpenSwapSuccess, sourceAccountID xdr.AccountId,
) ([]history2.ParticipantEffect, error) {

	sourceFixedFee := res.ActualSourceFee.Fixed
	sourcePercentFee := res.ActualSourceFee.Percent
	destFixedFee := res.ActualDestinationFee.Fixed
	destPercentFee := res.ActualDestinationFee.Percent
	if op.FeeData.SourcePaysForDest {
		sourceFixedFee += destFixedFee
		destFixedFee = 0
		sourcePercentFee += destPercentFee
		destPercentFee = 0
	}
	toLock := op.Amount + sourceFixedFee + sourcePercentFee

	source := h.BalanceEffect(op.SourceBalance, &history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(toLock),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(sourcePercentFee),
				Fixed:             regources.Amount(sourceFixedFee),
			},
		},
	})

	return []history2.ParticipantEffect{source}, nil
}
