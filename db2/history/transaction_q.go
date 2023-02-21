package history

import (
	"github.com/SafeRE-IT/horizon/db2"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/toid"
)

var selectTransaction = sq.Select(
	"ht.id, " +
		"ht.transaction_hash, " +
		"ht.ledger_sequence, " +
		"ht.application_order, " +
		"ht.account, " +
		"ht.fee_paid, " +
		"ht.operation_count, " +
		"ht.tx_envelope, " +
		"ht.tx_result, " +
		"ht.tx_meta, " +
		"ht.tx_fee_meta, " +
		"array_to_string(ht.signatures, ',') AS signatures, " +
		"ht.memo_type, " +
		"ht.memo, " +
		"lower(ht.time_bounds) AS valid_after, " +
		"upper(ht.time_bounds) AS valid_before, " +
		"ht.ledger_close_time").
	From("history_transactions ht")

// TransactionsQ is a helper struct to aid in configuring queries that loads
// slices of transaction structs.
type TransactionsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type TransactionsQI interface {
	ForAccount(aid string) TransactionsQI
	ForLedger(seq int32) TransactionsQI
	Page(page db2.PageQuery) TransactionsQI
	ByTxIDs(txIDs []int64) TransactionsQI
	Select(dest interface{}) error
}

// Transactions provides a helper to filter rows from the `history_transactions`
// table with pre-defined filters.  See `TransactionsQ` methods for the
// available filters.
func (q *Q) Transactions() TransactionsQI {
	return &TransactionsQ{
		parent: q,
		sql:    selectTransaction,
	}
}

// TransactionByHash is a query that loads a single row from the
// `history_transactions` table based upon the provided hash.
func (q *Q) TransactionByHash(dest interface{}, hash string) error {
	sql := selectTransaction.
		Limit(1).
		Where("ht.transaction_hash = ?", hash)

	return q.Get(dest, sql)
}

// TransactionByHashOrID is a query that loads a single row from the
// `history_transactions` table based upon the provided hash or id.
func (q *Q) TransactionByHashOrID(dest interface{}, hashOrID string) error {
	txID, err := strconv.ParseUint(hashOrID, 10, 64)
	if err != nil {
		return q.TransactionByHash(dest, hashOrID)
	}

	sql := selectTransaction.
		Limit(1).
		Where("ht.id = ?", txID)

	return q.Get(dest, sql)
}

// ForAccount filters the transactions collection to a specific account
func (q *TransactionsQ) ForAccount(aid string) TransactionsQI {
	var account Account
	q.Err = q.parent.AccountByAddress(&account, aid)
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.
		Join("history_transaction_participants htp ON htp.history_transaction_id = ht.id").
		Where("htp.history_account_id = ?", account.ID)

	return q
}

// ForLedger filters the query to a only transactions in a specific ledger,
// specified by its sequence.
func (q *TransactionsQ) ForLedger(seq int32) TransactionsQI {
	var ledger Ledger
	q.Err = q.parent.LedgerBySequence(&ledger, seq)
	if q.Err != nil {
		return q
	}

	start := toid.ID{LedgerSequence: seq}
	end := toid.ID{LedgerSequence: seq + 1}
	q.sql = q.sql.Where(
		"ht.id >= ? AND ht.id < ?",
		start.ToInt64(),
		end.ToInt64(),
	)

	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *TransactionsQ) Page(page db2.PageQuery) TransactionsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ht.id")
	return q
}

func (q *TransactionsQ) ByTxIDs(txIDs []int64) TransactionsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"ht.id": txIDs})
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *TransactionsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
