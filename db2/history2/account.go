package history2

// Account is a row of data from the `history_accounts` table
type Account struct {
	ID                uint64 `db:"id"`
	Address           string `db:"address"`
	KycRecoveryStatus int    `db:"kyc_recovery_status"`
}

//NewAccount - creates new instance of account
func NewAccount(accountID uint64, address string) Account {
	return Account{
		ID:      accountID,
		Address: address,
	}
}
