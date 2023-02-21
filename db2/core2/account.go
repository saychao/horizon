package core2

// Account is a row of data from the `accounts` table
type Account struct {
	Address    string  `db:"account_id"`
	Referrer   *string `db:"referrer"`
	SequenceID uint64  `db:"sequential_id"`
	RoleID     uint64  `db:"role_id"`
}
