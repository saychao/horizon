package history

import (
	"github.com/saychao/horizon/db2"
	"time"

	"github.com/guregu/null"
)

// Ledger is a row of data from the `history_ledgers` table
type Ledger struct {
	db2.TotalOrderID
	Sequence           int32       `db:"sequence"`
	ImporterVersion    int32       `db:"importer_version"`
	LedgerHash         string      `db:"ledger_hash"`
	PreviousLedgerHash null.String `db:"previous_ledger_hash"`
	TransactionCount   int32       `db:"transaction_count"`
	OperationCount     int32       `db:"operation_count"`
	ClosedAt           time.Time   `db:"closed_at"`
	TotalCoins         int64       `db:"total_coins"`
	FeePool            int64       `db:"fee_pool"`
	BaseFee            int32       `db:"base_fee"`
	BaseReserve        int32       `db:"base_reserve"`
	MaxTxSetSize       int32       `db:"max_tx_set_size"`
}
