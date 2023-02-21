package statslimits

import (
	"math"

	"github.com/saychao/horizon/db2/core2"
)

//Table is used to built complete limits & stats overview,
//using different level limits

var maxLimits = core2.Limits{
	DailyOut:   math.MaxUint64,
	WeeklyOut:  math.MaxUint64,
	MonthlyOut: math.MaxUint64,
	AnnualOut:  math.MaxUint64,
}

type (
	Group struct {
		AssetCode   string
		StatsOpType int32
	}
	LimitsWithStats struct {
		Limits core2.Limits
		Stats  core2.Statistics
	}
	Table map[Group]LimitsWithStats
)

func NewTable(limits []core2.Limits, stats []core2.Statistics) (lt Table) {
	lt = Table{}
	for _, entry := range limits {
		key := Group{
			AssetCode:   entry.AssetCode,
			StatsOpType: entry.StatsOpType,
		}

		limitsEntry := LimitsWithStats{
			Limits: entry,
		}

		lt[key] = limitsEntry
	}

	for _, entry := range stats {
		key := Group{
			AssetCode:   entry.AssetCode,
			StatsOpType: entry.OperationType,
		}

		limitsWithStatsEntry := LimitsWithStats{
			Limits: lt[key].Limits,
			Stats:  entry,
		}

		lt[key] = limitsWithStatsEntry
	}

	return lt
}

func (statslimitsTable Table) Update(limits []core2.Limits) {
	for _, v := range limits {
		key := Group{
			AssetCode:   v.AssetCode,
			StatsOpType: v.StatsOpType,
		}

		entry, ok := statslimitsTable[key]
		if !ok {
			entry = LimitsWithStats{
				Stats: core2.Statistics{},
			}
		}

		statslimitsTable[key] = LimitsWithStats{
			Limits: v,
			Stats:  entry.Stats,
		}
	}
}

func (statslimitsTable Table) FulfillEmptyLimits() {
	for k, v := range statslimitsTable {
		if v.Limits.ID != 0 {
			continue
		}
		statslimitsTable[k] = LimitsWithStats{
			Limits: maxLimits,
			Stats:  v.Stats,
		}
	}
}
