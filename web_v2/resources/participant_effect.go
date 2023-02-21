package resources

import (
	"strconv"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// NewParticipantEffectKey - returns new key for ID
func NewParticipantEffectKey(id int64) regources.Key {
	return regources.Key{
		Type: regources.PARTICIPANT_EFFECTS,
		ID:   strconv.FormatInt(id, 10),
	}
}

var participantEffectsType = map[regources.ResourceType]history2.EffectType{
	regources.EFFECTS_FUNDED:              history2.EffectTypeFunded,
	regources.EFFECTS_ISSUED:              history2.EffectTypeIssued,
	regources.EFFECTS_CHARGED:             history2.EffectTypeCharged,
	regources.EFFECTS_WITHDRAWN:           history2.EffectTypeWithdrawn,
	regources.EFFECTS_LOCKED:              history2.EffectTypeLocked,
	regources.EFFECTS_UNLOCKED:            history2.EffectTypeUnlocked,
	regources.EFFECTS_CHARGED_FROM_LOCKED: history2.EffectTypeChargedFromLocked,
	regources.EFFECTS_MATCHED:             history2.EffectTypeMatched,
}

func EffectTypeFromString(effectType regources.ResourceType) history2.EffectType {
	return participantEffectsType[effectType]
}

// NewEffect - returns new instance of effect
func NewEffect(id int64, effect history2.Effect) regources.Resource {
	switch effect.Type {
	case history2.EffectTypeFunded:
		return newBalanceChangeEffect(id, regources.EFFECTS_FUNDED, *effect.Funded)
	case history2.EffectTypeIssued:
		return newBalanceChangeEffect(id, regources.EFFECTS_ISSUED, *effect.Issued)
	case history2.EffectTypeCharged:
		return newBalanceChangeEffect(id, regources.EFFECTS_CHARGED, *effect.Charged)
	case history2.EffectTypeWithdrawn:
		return newBalanceChangeEffect(id, regources.EFFECTS_WITHDRAWN, *effect.Withdrawn)
	case history2.EffectTypeLocked:
		return newBalanceChangeEffect(id, regources.EFFECTS_LOCKED, *effect.Locked)
	case history2.EffectTypeUnlocked:
		return newBalanceChangeEffect(id, regources.EFFECTS_UNLOCKED, *effect.Unlocked)
	case history2.EffectTypeChargedFromLocked:
		return newBalanceChangeEffect(id, regources.EFFECTS_CHARGED_FROM_LOCKED, *effect.ChargedFromLocked)
	case history2.EffectTypeMatched:
		return newMatchedEffect(id, *effect.Matched)
	default:
		panic(errors.From(errors.New("unexpected effect type"), logan.F{
			"type": effect.Type,
		}))
	}
}

func newMatchedEffect(id int64, effect history2.MatchEffect) *regources.EffectMatched {
	return &regources.EffectMatched{
		Key: regources.Key{
			Type: regources.EFFECTS_MATCHED,
			ID:   strconv.FormatInt(id, 10),
		},
		Attributes: regources.EffectMatchedAttributes{
			OrderBookId: effect.OrderBookID,
			OfferId:     effect.OfferID,
			Price:       effect.Price,
			Charged:     newParticularBalanceChange(effect.Charged),
			Funded:      newParticularBalanceChange(effect.Funded),
		},
	}
}

func newParticularBalanceChange(effect history2.ParticularBalanceChangeEffect) regources.ParticularBalanceChangeEffect {
	return regources.ParticularBalanceChangeEffect{
		Amount:         effect.Amount,
		Fee:            effect.Fee,
		BalanceAddress: effect.BalanceAddress,
		AssetCode:      effect.AssetCode,
	}
}

func newBalanceChangeEffect(id int64, resourceType regources.ResourceType,
	effect history2.BalanceChangeEffect) *regources.EffectBalanceChange {
	return &regources.EffectBalanceChange{
		Key: regources.Key{
			Type: resourceType,
			ID:   strconv.FormatInt(id, 10),
		},
		Attributes: regources.EffectBalanceChangeAttributes{
			Amount: effect.Amount,
			Fee:    effect.Fee,
		},
	}
}
