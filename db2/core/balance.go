package core

type Balance struct {
	BalanceID string `db:"balance_id"`
	AccountID string `db:"account_id"`
	Asset     string `db:"asset"`
	Amount    int64  `db:"amount"`
	Locked    int64  `db:"locked"`
}

// DEPRECATED
func (q *Q) BalancesByAddress(dest interface{}, addy string) error {
	sql := selectBalance.Where("ba.account_id = ?", addy)
	return q.Select(dest, sql)
}

// DEPRECATED
func (q *Q) BalanceByID(dest interface{}, bid string) error {
	sql := selectBalance.Where("ba.balance_id = ?", bid)
	return q.Get(dest, sql)
}
