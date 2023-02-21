package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type cancelAtomicSwapAskOpHandler struct {
	effectsProvider
}

// Details returns details about cancel atomic swap bid operation
func (h *cancelAtomicSwapAskOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: xdr.OperationTypeCancelAtomicSwapAsk,
		CancelAtomicSwapAsk: &history2.CancelAtomicSwapAskDetails{
			BidID: int64(op.Body.MustCancelAtomicSwapAskOp().AskId),
		},
	}, nil
}

// ParticipantsEffects returns participants effects with source effect `unlocked`
// if atomic swap bid has zero locked amount
func (h *cancelAtomicSwapAskOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	opResult := opRes.MustCancelAtomicSwapAskResult().MustSuccess()

	// it means that there is pending atomic swap request,
	// so bid still exists
	// we must wait for review that request
	if opResult.LockedAmount != 0 {
		return h.effectsProvider.ParticipantsEffects(opBody, opRes, sourceAccountID, ledgerChanges)
	}

	atomicSwapAsk := h.getAtomicSwapAsk(opBody.MustCancelAtomicSwapAskOp().AskId, ledgerChanges)

	if atomicSwapAsk == nil {
		return nil, errors.From(
			errors.New("expected atomic swap to be in STATE ledger changes"), map[string]interface{}{
				"bid_id": uint64(opBody.MustCancelAtomicSwapAskOp().AskId),
			})
	}

	return []history2.ParticipantEffect{h.BalanceEffect(atomicSwapAsk.BaseBalance, &history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(atomicSwapAsk.Amount),
		},
	})}, nil
}

func (h *cancelAtomicSwapAskOpHandler) getAtomicSwapAsk(askID xdr.Uint64,
	ledgerChanges []xdr.LedgerEntryChange,
) *xdr.AtomicSwapAskEntry {
	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeAtomicSwapAsk {
			continue
		}

		atomicSwapBid := change.MustState().Data.MustAtomicSwapAsk()

		if atomicSwapBid.Id == askID {
			return &atomicSwapBid
		}
	}

	return nil
}
