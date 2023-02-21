package history2

// ParticipantEffect - stores effect of operation done to particular account or balance
type ParticipantEffect struct {
	ID          int64   `db:"id"`
	AccountID   uint64  `db:"account_id"`
	BalanceID   *uint64 `db:"balance_id"`
	AssetCode   *string `db:"asset_code"`
	Effect      *Effect `db:"effect"`
	OperationID int64   `db:"operation_id"`
	// Operation - is populated on join 1 to 1 relation is quarantined
	*Operation `db:"operations"`
	//BalanceAddress - is populate on join of balance. 1 to [0:1] relation
	BalanceAddress *string `db:"balance_address"`
	//AccountAddress - is populate on join of account. 1 to 1 relation
	AccountAddress string `db:"account_address"`
}
