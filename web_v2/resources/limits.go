package resources

import (
	"github.com/saychao/horizon/db2/core2"
	"github.com/spf13/cast"
	regources "gitlab.com/tokend/regources/generated"
)

// NewLimits creates new instance of Limits from provided one
func NewLimits(limits core2.Limits) regources.Limits {
	newLimits := regources.Limits{
		Key: NewLimitsKey(limits.ID),
		Attributes: regources.LimitsAttributes{
			StatsOpType:     limits.StatsOpType,
			IsConvertNeeded: limits.IsConvertNeeded,
			DailyOut:        regources.Amount(limits.DailyOut),
			WeeklyOut:       regources.Amount(limits.WeeklyOut),
			MonthlyOut:      regources.Amount(limits.MonthlyOut),
			AnnualOut:       regources.Amount(limits.AnnualOut),
		},
		Relationships: regources.LimitsRelationships{
			Asset: NewAssetKey(limits.AssetCode).AsRelation(),
		},
	}

	if limits.AccountId != nil {
		newLimits.Relationships.Account = NewAccountKey(*limits.AccountId).AsRelation()
	}

	if limits.AccountType != nil {
		newLimits.Relationships.AccountRole = NewAccountRoleKey(*limits.AccountType).AsRelation()
	}

	return newLimits
}

func NewLimitsKey(id uint64) regources.Key {
	return regources.Key{
		ID:   cast.ToString(id),
		Type: regources.LIMITS,
	}
}
