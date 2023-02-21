package history2

//Balance - represents balance instance used to optimize indexing of db
type Balance struct {
	ID        uint64 `db:"id"`
	AccountID uint64 `db:"account_id"`
	Address   string `db:"address"`
	AssetCode string `db:"asset_code"`
}

//NewBalance - creates new instance of Balance
func NewBalance(balanceID, accountID uint64, address, assetCode string) Balance {
	return Balance{
		ID:        balanceID,
		AccountID: accountID,
		Address:   address,
		AssetCode: assetCode,
	}
}
