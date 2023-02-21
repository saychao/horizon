package history

type Balance struct {
	ID        int64  `db:"id"`
	BalanceID string `db:"balance_id"`
	AccountID string `db:"account_id"`
	Asset     string `db:"asset"`
}
