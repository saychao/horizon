package core2

type Fee struct {
	FeeType     int32  `db:"fee_type"`
	Asset       string `db:"asset"`
	Fixed       int64  `db:"fixed"`
	Percent     int64  `db:"percent"`
	Subtype     int64  `db:"subtype"`
	AccountID   string `db:"account_id"`
	AccountRole uint64 `db:"account_role"`
	LowerBound  int64  `db:"lower_bound"`
	UpperBound  int64  `db:"upper_bound"`
	Hash        string `db:"hash"`

	LastModified int32 `db:"lastmodified"`
}
