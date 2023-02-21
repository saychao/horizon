package resource

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
)

type StatisticsV2 struct {
	Id              uint64 `json:"id"`
	AccountId       string `json:"account_id"`
	StatsOpType     int32  `json:"stats_op_type"`
	AssetCode       string `json:"asset_code"`
	IsConvertNeeded bool   `json:"is_convert_needed"`
	DailyOutcome    string `json:"daily_outcome"`
	WeeklyOutcome   string `json:"weekly_outcome"`
	MonthlyOutcome  string `json:"monthly_outcome"`
	AnnualOutcome   string `json:"annual_outcome"`
	UpdatedAt       int64  `json:"updated_at"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *StatisticsV2) Populate(row core.StatisticsV2Entry) {
	s.Id = row.Id
	s.AccountId = row.AccountId
	s.StatsOpType = row.StatsOpType
	s.AssetCode = row.AssetCode
	s.IsConvertNeeded = row.IsConvertNeeded
	s.DailyOutcome = amount.StringU(row.DailyOutcome)
	s.WeeklyOutcome = amount.StringU(row.WeeklyOutcome)
	s.MonthlyOutcome = amount.StringU(row.MonthlyOutcome)
	s.AnnualOutcome = amount.StringU(row.AnnualOutcome)
	s.UpdatedAt = row.UpdatedAt
}
