package limits

import (
	"github.com/saychao/horizon/db2/core2"
)

//Table is used to built complete limits overview,
//using different level limits

type (
	Group struct {
		AssetCode   string
		StatsOpType int32
	}
	Table map[Group]core2.Limits
)

func NewTable(limits []core2.Limits) (lt Table) {
	lt = Table{}
	for _, entry := range limits {
		key := Group{
			AssetCode:   entry.AssetCode,
			StatsOpType: entry.StatsOpType,
		}

		lt[key] = entry
	}

	return lt
}

func (lt Table) Update(limits []core2.Limits) {
	for _, v := range limits {
		key := Group{
			AssetCode:   v.AssetCode,
			StatsOpType: v.StatsOpType,
		}

		lt[key] = v
	}
}
