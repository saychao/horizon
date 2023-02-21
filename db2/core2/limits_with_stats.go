package core2

type LimitsWithStats struct {
	ID          string
	AccountID   string
	AssetCode   string
	StatsOpType int32
	Statistics  Statistics
	Limits      Limits
}
