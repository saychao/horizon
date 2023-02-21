package statslimits

import (
	"fmt"
	"github.com/SafeRE-IT/horizon/db2/core2"
)

func (statslimitsTable Table) CoreUnitsList() []core2.LimitsWithStats {
	res := make([]core2.LimitsWithStats, 0, len(statslimitsTable))

	for g, coreUnit := range statslimitsTable {
		res = append(res, core2.LimitsWithStats{
			ID:          fmt.Sprintf("%d:%d", coreUnit.Limits.ID, coreUnit.Stats.ID),
			AccountID:   coreUnit.Stats.AccountID,
			StatsOpType: g.StatsOpType,
			AssetCode:   g.AssetCode,
			Limits:      coreUnit.Limits,
			Statistics:  coreUnit.Stats,
		})

	}

	return res
}
