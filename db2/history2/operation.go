package history2

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

// Operation - stores details of operation performed
type Operation struct {
	ID              int64             `db:"id"`
	TxID            int64             `db:"tx_id"`
	Type            xdr.OperationType `db:"type"`
	Details         OperationDetails  `db:"details"`
	LedgerCloseTime time.Time         `db:"ledger_close_time"`
	Source          string            `db:"source"`
}
