package history

import (
	"github.com/saychao/horizon/db2"
	"time"

	"github.com/guregu/null"
)

// Transaction is a row of data from the `history_transactions` table
type Transaction struct {
	db2.TotalOrderID
	TransactionHash  string      `db:"transaction_hash"`
	LedgerSequence   int32       `db:"ledger_sequence"`
	LedgerCloseTime  time.Time   `db:"ledger_close_time"`
	ApplicationOrder int32       `db:"application_order"`
	Account          string      `db:"account"`
	FeePaid          int32       `db:"fee_paid"`
	OperationCount   int32       `db:"operation_count"`
	TxEnvelope       string      `db:"tx_envelope"`
	TxResult         string      `db:"tx_result"`
	TxMeta           string      `db:"tx_meta"`
	TxFeeMeta        string      `db:"tx_fee_meta"`
	SignatureString  string      `db:"signatures"`
	MemoType         string      `db:"memo_type"`
	Memo             null.String `db:"memo"`
	ValidAfter       null.Int    `db:"valid_after"`
	ValidBefore      null.Int    `db:"valid_before"`
}
