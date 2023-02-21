package core

// Account is a row of data from the `accounts` table
type Account struct {
	AccountID    string  `db:"account_id"`
	Referrer     *string `db:"referrer"`
	SequentialID uint64  `db:"sequential_id"`
	RoleID       uint64  `db:"role_id"`
	StatisticsV2 []StatisticsV2Entry
	*AccountKYC
}
