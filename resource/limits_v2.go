package resource

import (
	"github.com/saychao/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
)

type LimitsV2 struct {
	Id              uint64  `json:"id"`
	AccountType     *int32  `json:"account_type"`
	AccountId       *string `json:"account_id"`
	StatsOpType     int32   `json:"stats_op_type"`
	AssetCode       string  `json:"asset_code"`
	IsConvertNeeded bool    `json:"is_convert_needed"`
	DailyOut        string  `json:"daily_out"`
	WeeklyOut       string  `json:"weekly_out"`
	MonthlyOut      string  `json:"monthly_out"`
	AnnualOut       string  `json:"annual_out"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *LimitsV2) Populate(row core.LimitsV2Entry) {
	s.Id = row.Id
	s.AccountType = row.AccountType
	s.AccountId = row.AccountId
	s.StatsOpType = row.StatsOpType
	s.AssetCode = row.AssetCode
	s.IsConvertNeeded = row.IsConvertNeeded
	s.DailyOut = amount.StringU(row.DailyOut)
	s.WeeklyOut = amount.StringU(row.WeeklyOut)
	s.MonthlyOut = amount.StringU(row.MonthlyOut)
	s.AnnualOut = amount.StringU(row.AnnualOut)
}
