package core

import (
	"math"
)

type LimitsV2Entry struct {
	Id              uint64  `db:"id"`
	AccountType     *int32  `db:"account_type"`
	AccountId       *string `db:"account_id"`
	StatsOpType     int32   `db:"stats_op_type"`
	AssetCode       string  `db:"asset_code"`
	IsConvertNeeded bool    `db:"is_convert_needed"`
	DailyOut        uint64  `db:"daily_out"`
	WeeklyOut       uint64  `db:"weekly_out"`
	MonthlyOut      uint64  `db:"monthly_out"`
	AnnualOut       uint64  `db:"annual_out"`
}

func DefaultLimits(accountType *int32, accountID *string, statsOpType int32, assetCode string) LimitsV2Entry {
	return LimitsV2Entry{
		AccountType: accountType,
		AccountId:   accountID,
		StatsOpType: statsOpType,
		AssetCode:   assetCode,
		DailyOut:    math.MaxInt64,
		WeeklyOut:   math.MaxInt64,
		MonthlyOut:  math.MaxInt64,
		AnnualOut:   math.MaxInt64,
	}

}
