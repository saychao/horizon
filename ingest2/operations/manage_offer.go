package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type manageOfferOpHandler struct {
	effectsProvider
}

// Details returns details about manage offer operation
func (h *manageOfferOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageOfferOp := op.Body.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	offerID := int64(manageOfferOp.OfferId)
	isDeleted := manageOfferOpRes.Offer.Effect == xdr.ManageOfferEffectDeleted
	if !isDeleted {
		offerID = int64(manageOfferOpRes.Offer.MustOffer().OfferId)
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageOffer,
		ManageOffer: &history2.ManageOfferDetails{
			OfferID:     offerID,
			OrderBookID: int64(manageOfferOp.OrderBookId),
			BaseAsset:   string(manageOfferOpRes.BaseAsset),
			QuoteAsset:  string(manageOfferOpRes.QuoteAsset),
			Amount:      regources.Amount(manageOfferOp.Amount),
			Price:       regources.Amount(manageOfferOp.Price),
			IsBuy:       manageOfferOp.IsBuy,
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(manageOfferOp.Fee),
			},
			IsDeleted: isDeleted,
		},
	}, nil
}

// ParticipantsEffects can return `matched` and `locked` effects if offer created
// returns `unlocked` effects if offer canceled (deleted by user)
func (h *manageOfferOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageOfferOp := opBody.MustManageOfferOp()
	manageOfferOpRes := opRes.MustManageOfferResult().MustSuccess()

	if manageOfferOp.Amount != 0 {
		source := h.Participant(sourceAccountID)
		return h.getNewOfferEffect(manageOfferOp, manageOfferOpRes, source, ledgerChanges), nil
	}

	deletedOfferEffects := h.getDeletedOffersEffect(ledgerChanges)
	if len(deletedOfferEffects) != 1 {
		return nil, errors.From(errors.New("Unexpected number of deleted offer for manage offer delete"), logan.F{
			"expected": 1,
			"actual":   len(deletedOfferEffects),
		})
	}

	return deletedOfferEffects, nil
}

func (h *manageOfferOpHandler) getNewOfferEffect(op xdr.ManageOfferOp,
	res xdr.ManageOfferSuccessResult, source history2.ParticipantEffect, ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {
	participants := make([]history2.ParticipantEffect, 0, 4)

	if op.OrderBookId != 0 {
		sale := getImmediateSale(ledgerChanges, op.OrderBookId)
		if sale != nil {
			participants = append(participants, h.BalanceEffect(sale.BaseBalance, &history2.Effect{
				Type: history2.EffectTypeIssued,
				Issued: &history2.BalanceChangeEffect{
					Amount: regources.Amount(op.Amount),
					Fee:    regources.Fee{},
				},
			}), h.BalanceEffect(sale.BaseBalance, &history2.Effect{
				Type: history2.EffectTypeLocked,
				Locked: &history2.BalanceChangeEffect{
					Amount: regources.Amount(op.Amount),
					// we do not fee because we locked base (fee in quote)
					Fee: regources.Fee{},
				},
			}))
		}
	}

	matchedParticipants, _ := h.getMatchesEffects(res.OffersClaimed, offer{
		OrderBookID:         int64(op.OrderBookId),
		AccountID:           source.AccountID,
		BaseBalanceID:       h.MustBalanceID(op.BaseBalance),
		BaseBalanceAddress:  op.BaseBalance.AsString(),
		QuoteBalanceID:      h.MustBalanceID(op.QuoteBalance),
		QuoteBalanceAddress: op.QuoteBalance.AsString(),
		BaseAsset:           string(res.BaseAsset),
		QuoteAsset:          string(res.QuoteAsset),
		IsBuy:               op.IsBuy,
	})
	participants = append(participants, matchedParticipants...)
	if res.Offer.Effect == xdr.ManageOfferEffectDeleted {
		return participants
	}

	// we need to handle amount which was not matched
	newOffer := res.Offer.MustOffer()
	source.AssetCode = new(string)
	source.BalanceID = new(uint64)
	source.Effect = &history2.Effect{
		Type:   history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{},
	}
	if newOffer.IsBuy {
		*source.BalanceID = h.MustBalanceID(newOffer.QuoteBalance)
		*source.AssetCode = string(newOffer.Quote)
		source.Effect.Locked.Amount = regources.Amount(newOffer.QuoteAmount)
		source.Effect.Locked.Fee.CalculatedPercent = regources.Amount(newOffer.Fee)
	} else {
		*source.BalanceID = h.MustBalanceID(newOffer.BaseBalance)
		*source.AssetCode = string(newOffer.Base)
		source.Effect.Locked.Amount = regources.Amount(newOffer.BaseAmount)
	}

	participants = append(participants, source)
	return participants
}

type offerID struct {
	OrderBookID uint64
	OfferID     uint64
}

// getDeletedOffersEffect - creates participant effects for all offer entries with change type `State`
func (h *manageOfferOpHandler) getDeletedOffersEffect(ledgerChanges []xdr.LedgerEntryChange,
) []history2.ParticipantEffect {

	result := map[offerID]history2.ParticipantEffect{}
	for _, change := range ledgerChanges {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeOfferEntry {
			continue
		}

		deletedOffer := change.MustState().Data.MustOffer()

		participant := history2.ParticipantEffect{
			AccountID: h.MustAccountID(deletedOffer.OwnerId),
			BalanceID: new(uint64),
			AssetCode: new(string),
			Effect: &history2.Effect{
				Type:     history2.EffectTypeUnlocked,
				Unlocked: &history2.BalanceChangeEffect{},
			},
		}

		if deletedOffer.IsBuy {
			*participant.BalanceID = h.MustBalanceID(deletedOffer.QuoteBalance)
			*participant.AssetCode = string(deletedOffer.Quote)
			participant.Effect.Unlocked.Amount = regources.Amount(deletedOffer.QuoteAmount)
			participant.Effect.Unlocked.Fee.CalculatedPercent = regources.Amount(deletedOffer.Fee)
		} else {
			*participant.BalanceID = h.MustBalanceID(deletedOffer.BaseBalance)
			*participant.AssetCode = string(deletedOffer.Base)
			participant.Effect.Unlocked.Amount = regources.Amount(deletedOffer.BaseAmount)
		}

		result[offerID{
			OfferID:     uint64(deletedOffer.OfferId),
			OrderBookID: uint64(deletedOffer.OrderBookId),
		}] = participant
	}

	uniqueResults := make([]history2.ParticipantEffect, 0, len(result))
	for _, participant := range result {
		uniqueResults = append(uniqueResults, participant)
	}

	return uniqueResults
}

type offer struct {
	OrderBookID         int64
	AccountID           uint64
	BaseBalanceID       uint64
	BaseBalanceAddress  string
	QuoteBalanceID      uint64
	QuoteBalanceAddress string
	BaseAsset           string
	QuoteAsset          string
	IsBuy               bool
}

// getParticipantsEffects - returns participants effects based on the provided matches and total base amount
func (h *manageOfferOpHandler) getMatchesEffects(claimOfferAtoms []xdr.ClaimOfferAtom,
	sourceOffer offer) ([]history2.ParticipantEffect, uint64) {
	var totalBaseAmount uint64
	result := make([]history2.ParticipantEffect, 0, len(claimOfferAtoms)*4)

	priceLevels := make(map[xdr.Int64][]xdr.ClaimOfferAtom)
	for _, matchedOffer := range claimOfferAtoms {
		priceLevels[matchedOffer.CurrentPrice] = append(priceLevels[matchedOffer.CurrentPrice], matchedOffer)
	}

	for price, matches := range priceLevels {
		var levelBaseAmount xdr.Int64 = 0
		var levelQuoteAmount xdr.Int64 = 0
		var levelFee xdr.Int64 = 0
		for _, match := range matches {
			result = h.addParticipantEffects(result, offer{
				OrderBookID:         sourceOffer.OrderBookID,
				AccountID:           h.MustAccountID(match.BAccountId),
				BaseBalanceID:       h.MustBalanceID(match.BaseBalance),
				QuoteBalanceID:      h.MustBalanceID(match.QuoteBalance),
				BaseBalanceAddress:  match.BaseBalance.AsString(),
				QuoteBalanceAddress: match.QuoteBalance.AsString(),
				BaseAsset:           sourceOffer.BaseAsset,
				QuoteAsset:          sourceOffer.QuoteAsset,
				IsBuy:               !sourceOffer.IsBuy,
			}, int64(match.OfferId), match.BaseAmount, match.QuoteAmount, match.CurrentPrice, match.BFeePaid, true)

			levelBaseAmount += match.BaseAmount
			levelQuoteAmount += match.QuoteAmount
			levelFee += match.AFeePaid
		}

		result = h.addParticipantEffects(result, sourceOffer, 0, levelBaseAmount,
			levelQuoteAmount, price, levelFee, false)

		totalBaseAmount += uint64(levelBaseAmount)
	}

	return result, totalBaseAmount
}

func (h *manageOfferOpHandler) addParticipantEffects(participants []history2.ParticipantEffect,
	offer offer, id int64, baseAmount, quoteAmount, price, fee xdr.Int64, IsUnlocked bool,
) []history2.ParticipantEffect {
	baseBalanceEffect := history2.ParticularBalanceChangeEffect{
		BalanceAddress: offer.BaseBalanceAddress,
		AssetCode:      offer.BaseAsset,
		BalanceChangeEffect: history2.BalanceChangeEffect{
			Amount: regources.Amount(baseAmount),
		},
	}

	quoteBalanceEffect := history2.ParticularBalanceChangeEffect{
		BalanceAddress: offer.QuoteBalanceAddress,
		AssetCode:      offer.QuoteAsset,
		BalanceChangeEffect: history2.BalanceChangeEffect{
			Amount: regources.Amount(quoteAmount),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(fee),
			},
		},
	}

	matchedOfferEffect := history2.Effect{
		Type: history2.EffectTypeMatched,
		Matched: &history2.MatchEffect{
			OfferID:     id,
			OrderBookID: offer.OrderBookID,
			Price:       regources.Amount(price),
		},
	}

	if offer.IsBuy {
		matchedOfferEffect.Matched.Funded = baseBalanceEffect
		matchedOfferEffect.Matched.Charged = quoteBalanceEffect
		if IsUnlocked {
			participants = append(participants, setUnlocked(offer, quoteBalanceEffect, true))
		}
	} else {
		matchedOfferEffect.Matched.Funded = quoteBalanceEffect
		matchedOfferEffect.Matched.Charged = baseBalanceEffect
		if IsUnlocked {
			participants = append(participants, setUnlocked(offer, baseBalanceEffect, false))
		}
	}

	return append(participants,
		history2.ParticipantEffect{
			AccountID: offer.AccountID,
			BalanceID: &offer.BaseBalanceID,
			AssetCode: &offer.BaseAsset,
			Effect:    &matchedOfferEffect,
		}, history2.ParticipantEffect{
			AccountID: offer.AccountID,
			BalanceID: &offer.QuoteBalanceID,
			AssetCode: &offer.QuoteAsset,
			Effect:    &matchedOfferEffect,
		})
}

func setUnlocked(offer offer, balanceEffect history2.ParticularBalanceChangeEffect,
	IsQuote bool,
) history2.ParticipantEffect {
	unlockedEffect := history2.Effect{
		Type: history2.EffectTypeUnlocked,
		Unlocked: &history2.BalanceChangeEffect{
			Amount: balanceEffect.Amount,
			Fee:    balanceEffect.Fee,
		},
	}

	if IsQuote {
		return history2.ParticipantEffect{
			AccountID: offer.AccountID,
			BalanceID: &offer.QuoteBalanceID,
			AssetCode: &offer.QuoteAsset,
			Effect:    &unlockedEffect}
	}

	return history2.ParticipantEffect{
		AccountID: offer.AccountID,
		BalanceID: &offer.BaseBalanceID,
		AssetCode: &offer.BaseAsset,
		Effect:    &unlockedEffect,
	}

}

func getImmediateSale(changes []xdr.LedgerEntryChange, saleID xdr.Uint64) *xdr.SaleEntry {
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryChangeTypeState {
			continue
		}

		if change.MustState().Data.Type != xdr.LedgerEntryTypeSale {
			continue
		}

		sale := change.MustState().Data.MustSale()
		if sale.SaleId != saleID {
			continue
		}

		if sale.SaleTypeExt.SaleType != xdr.SaleTypeImmediate {
			return nil
		}

		return &sale
	}

	return nil
}
