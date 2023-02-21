package core2

type Statistics struct {
	ID              uint64 `db:"id"`
	AccountID       string `db:"account_id"`
	OperationType   int32  `db:"stats_op_type"`
	AssetCode       string `db:"asset_code"`
	IsConvertNeeded bool   `db:"is_convert_needed"`
	DailyOutcome    uint64 `db:"daily_out"`
	WeeklyOutcome   uint64 `db:"weekly_out"`
	MonthlyOutcome  uint64 `db:"monthly_out"`
	AnnualOutcome   uint64 `db:"annual_out"`
	UpdatedAt       int64  `db:"updated_at"`
}
