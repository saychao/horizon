package history2

import (
	"github.com/saychao/horizon/db2"
	"time"
)

// Transaction is a row of data from the `history_transactions` table
type Transaction struct {
	db2.TotalOrderID
	Hash             string    `db:"hash"`
	LedgerSequence   int32     `db:"ledger_sequence"`
	LedgerCloseTime  time.Time `db:"ledger_close_time"`
	ApplicationOrder int32     `db:"application_order"`
	Account          string    `db:"account"`
	OperationCount   int32     `db:"operation_count"`
	Envelope         string    `db:"envelope"`
	Result           string    `db:"result"`
	Meta             string    `db:"meta"`
	ValidAfter       time.Time `db:"valid_after"`
	ValidBefore      time.Time `db:"valid_before"`
}
