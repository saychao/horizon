package resources

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewLimitsWithStats(record *core2.LimitsWithStats) regources.LimitsWithStats {
	return regources.LimitsWithStats{
		Key: regources.Key{
			ID:   record.ID,
			Type: regources.LIMITS_WITH_STATS,
		},
		Attributes: regources.LimitsWithStatsAttributes{
			Statistics: NewStatistics(record.Statistics),
			Limits:     NewLimits(record.Limits),
		},
		Relationships: regources.LimitsWithStatsRelationships{
			Account: NewAccountKey(record.AccountID).AsRelation(),
		},
	}
}
