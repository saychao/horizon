package history

// Account is a row of data from the `history_accounts` table
type Account struct {
	ID          int64  `db:"id"`
	Address     string `db:"address"`
	AccountType int32  `db:"account_type"`
}
