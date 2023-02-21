package operations

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

type checkSaleStateOpHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// Details returns details about check sale state operation
func (h *checkSaleStateOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: xdr.OperationTypeCheckSaleState,
		CheckSaleState: &history2.CheckSaleStateDetails{
			SaleID: int64(op.Body.MustCheckSaleStateOp().SaleId),
			Effect: opRes.MustCheckSaleStateResult().MustSuccess().Effect.Effect,
		},
	}, nil
}

// ParticipantsEffects returns sale owner and participants `matched` effects if sale closed
// returns `unlocked` effects if sale canceled or updated
func (h *checkSaleStateOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	res := opRes.MustCheckSaleStateResult().MustSuccess()

	switch res.Effect.Effect {
	case xdr.CheckSaleStateEffectCanceled, xdr.CheckSaleStateEffectUpdated:
		return h.manageOfferOpHandler.getDeletedOffersEffect(ledgerChanges), nil
	case xdr.CheckSaleStateEffectClosed:
		return h.getParticipationChanges(int64(opBody.MustCheckSaleStateOp().SaleId), res.Effect.MustSaleClosed(), ledgerChanges), nil
	default:
		return nil, errors.From(errors.New("unexpected check sale state result effect"), map[string]interface{}{
			"effect_i": int32(res.Effect.Effect),
			"sale_id":  uint64(res.SaleId),
		})
	}
}

func (h *checkSaleStateOpHandler) getParticipationChanges(orderBookID int64, closedRes xdr.CheckSaleClosedResult,
	ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {
	// TODO: we are not handling here cases that some parts of offers might be canceled due to rounding
	if len(closedRes.Results) == 0 {
		return nil
	}

	result := make([]history2.ParticipantEffect, 0, len(closedRes.Results)+1)
	var totalBaseIssued uint64
	ownerID := h.manageOfferOpHandler.MustAccountID(closedRes.SaleOwner)
	// it does not matter which base balance is used as we are sure that the operation of distribution will be clean
	baseBalanceAddress := closedRes.Results[0].SaleBaseBalance.AsString()
	baseBalanceID := h.manageOfferOpHandler.MustBalanceID(closedRes.Results[0].SaleBaseBalance)
	baseAsset := string(closedRes.Results[0].SaleDetails.BaseAsset)

	// we need to show explicitly that issuance has been perform to ensure that balance history is consistent
	// amount will be updated later
	result = append(result, history2.ParticipantEffect{
		AccountID: ownerID,
		BalanceID: &baseBalanceID,
		AssetCode: &baseAsset,
		Effect: &history2.Effect{
			Type: history2.EffectTypeIssued,
			Issued: &history2.BalanceChangeEffect{
				Amount: regources.Amount(0),
			},
		},
	})

	matchedOfferIDs := make(map[int64]struct{}, len(closedRes.Results)+1)

	for _, assetPairResult := range closedRes.Results {
		sourceOffer := offer{
			OrderBookID:         orderBookID,
			AccountID:           ownerID,
			BaseBalanceID:       baseBalanceID,
			BaseBalanceAddress:  baseBalanceAddress,
			QuoteBalanceID:      h.manageOfferOpHandler.MustBalanceID(assetPairResult.SaleQuoteBalance),
			QuoteBalanceAddress: assetPairResult.SaleQuoteBalance.AsString(),
			BaseAsset:           baseAsset,
			QuoteAsset:          string(assetPairResult.SaleDetails.QuoteAsset),
			IsBuy:               false,
		}
		assetPairMatches, baseIssued := h.manageOfferOpHandler.getMatchesEffects(
			assetPairResult.SaleDetails.OffersClaimed, sourceOffer)

		totalBaseIssued += baseIssued
		result = append(result, assetPairMatches...)

		for _, claimedOffer := range assetPairResult.SaleDetails.OffersClaimed {
			matchedOfferIDs[int64(claimedOffer.OfferId)] = struct{}{}
		}
	}

	// participant effects for offers not matched on sale close
	removedOffers := h.getRemovedOfferEntries(ledgerChanges, matchedOfferIDs)
	removedOfferEffects := h.getUnlockedEffects(removedOffers)
	result = append(result, removedOfferEffects...)

	result[0].Effect.Issued.Amount = regources.Amount(totalBaseIssued)

	return result
}

func (h *checkSaleStateOpHandler) getRemovedOfferEntries(changes []xdr.LedgerEntryChange, matchedOfferIDs map[int64]struct{},
) map[int64]xdr.OfferEntry {
	statedOffers := make(map[int64]xdr.OfferEntry)
	result := make(map[int64]xdr.OfferEntry)
	removedOffers := make(map[int64]bool)

	for _, change := range changes {
		switch change.Type {
		case xdr.LedgerEntryChangeTypeRemoved:
			if change.Removed.Type != xdr.LedgerEntryTypeOfferEntry {
				continue
			}

			offerID := int64(change.Removed.MustOffer().OfferId)

			if _, ok := matchedOfferIDs[offerID]; ok {
				continue
			}

			removedOffers[offerID] = true
		case xdr.LedgerEntryChangeTypeState:
			if change.State.Data.Type != xdr.LedgerEntryTypeOfferEntry {
				continue
			}

			offerID := int64(change.State.Data.MustOffer().OfferId)

			if _, ok := matchedOfferIDs[offerID]; ok {
				continue
			}

			statedOffers[offerID] = change.State.Data.MustOffer()
		}
	}

	for id := range removedOffers {
		if offer, ok := statedOffers[id]; ok {
			result[id] = offer
		}
	}

	return result
}

func (h *checkSaleStateOpHandler) getUnlockedEffects(removedOffers map[int64]xdr.OfferEntry) []history2.ParticipantEffect {
	result := map[offerID]history2.ParticipantEffect{}
	for _, off := range removedOffers {
		balanceID := off.BaseBalance
		if off.IsBuy {
			balanceID = off.QuoteBalance
		}
		unlockedAmount := off.BaseAmount
		if off.IsBuy {
			unlockedAmount = off.QuoteAmount
		}

		participant := h.manageOfferOpHandler.BalanceEffect(balanceID, &history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(unlockedAmount),
			},
		})

		if off.IsBuy {
			participant.Effect.Unlocked.Fee.CalculatedPercent = regources.Amount(off.Fee)
		}

		result[offerID{
			OrderBookID: uint64(off.OrderBookId),
			OfferID:     uint64(off.OfferId),
		}] = participant
	}

	uniqueResults := make([]history2.ParticipantEffect, 0, len(result))
	for _, participant := range result {
		uniqueResults = append(uniqueResults, participant)
	}

	return uniqueResults
}
