package history2

import "gitlab.com/tokend/go/xdr"

// LedgerChanges is a row of data from the `ledgers_changes` table
type LedgerChanges struct {
	TransactionID int64                 `db:"tx_id"`
	OperationID   int64                 `db:"op_id"`
	OrderNumber   int                   `db:"order_number"`
	Effect        int                   `db:"effect"`
	EntryType     int                   `db:"entry_type"`
	Payload       xdr.LedgerEntryChange `db:"payload"`
}
