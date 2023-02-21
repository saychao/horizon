package history2

import (
	"github.com/SafeRE-IT/horizon/db2"
	"time"
)

// Ledger is a row of data from the `ledgers` table
type Ledger struct {
	db2.TotalOrderID
	Sequence     int32     `db:"sequence"`
	Hash         string    `db:"hash"`
	PreviousHash string    `db:"previous_hash"`
	ClosedAt     time.Time `db:"closed_at"`
	TxCount      int32     `db:"tx_count"`
	Data         string    `db:"data"`
}
