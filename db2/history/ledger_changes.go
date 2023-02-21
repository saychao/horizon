package history

// LedgerChanges is a row of data from the `history_ledgers_changes` table
type LedgerChanges struct {
	TransactionID int64  `db:"tx_id"`
	OperationID   int64  `db:"op_id"`
	OrderNumber   int    `db:"order_number"`
	Effect        int    `db:"effect"`
	EntryType     int    `db:"entry_type"`
	Payload       string `db:"payload"`
}
