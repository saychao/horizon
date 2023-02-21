package history

import (
	"fmt"
	"github.com/SafeRE-IT/horizon/db2"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var selectLedgerChanges = sq.Select(
	"hlc.tx_id",
	"hlc.op_id",
	"hlc.order_number",
	"hlc.effect",
	"hlc.entry_type",
	"hlc.payload",
).From("history_ledger_changes hlc")

// LedgerChangesQ is a helper struct to aid in configuring queries that loads
// slices of LedgerChanges structs.
type LedgerChangesQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type LedgerChangesQI interface {
	// ByEffects filters query by specific effects
	ByEffects(effects []int) LedgerChangesQI
	// ByEntryType filters query by specific entry types
	ByEntryType(entryType []int) LedgerChangesQI
	// ByTransactionIDs filters query by specific tx_ids which based on page params,
	// entry types, effects
	ByTransactionIDs(page db2.PageQuery, entryTypes []int, effects []int) LedgerChangesQI
	// Select loads the results of the query specified by `q` into `dest`.
	Select(dest interface{}) error
}

// LedgerChanges provides a helper to filter rows from the `history_ledger_changes`
// table with pre-defined filters. See `LedgerChangesQ` methods for the available filters.
func (q *Q) LedgerChanges() LedgerChangesQI {
	return &LedgerChangesQ{
		parent: q,
		sql:    selectLedgerChanges,
	}
}

// ByEffects filters query by specific effects, if array empty do nothing
func (q *LedgerChangesQ) ByEffects(effects []int) LedgerChangesQI {
	if q.Err != nil {
		return q
	}

	if len(effects) == 0 {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"effect": effects})
	return q
}

// ByEntryType filters query by specific entry types, if array empty do nothing
func (q *LedgerChangesQ) ByEntryType(entryTypes []int) LedgerChangesQI {
	if q.Err != nil {
		return q
	}

	if len(entryTypes) == 0 {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"entry_type": entryTypes})
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *LedgerChangesQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

// ByTransactionsIDs filters ledger changes by tx_id which specifies on
// entry types, effects and paging params
func (q *LedgerChangesQ) ByTransactionIDs(page db2.PageQuery, entryTypes []int, effects []int) LedgerChangesQI {
	if q.Err != nil {
		return q
	}

	selectTxIDs := sq.Select("hlc.tx_id").From("history_ledger_changes hlc")
	if len(effects) != 0 {
		selectTxIDs = selectTxIDs.Where(sq.Eq{"effect": effects})
	}
	if len(entryTypes) != 0 {
		selectTxIDs = selectTxIDs.Where(sq.Eq{"entry_type": entryTypes})
	}

	selectTxIDs, q.Err = page.ApplyTo(selectTxIDs, "hlc.tx_id")

	if q.Err != nil {
		q.Err = errors.Wrap(q.Err, "failed to get paging params")
		return q
	}

	var sqlTxIDs string
	var args []interface{}
	sqlTxIDs, args, q.Err = selectTxIDs.ToSql()
	if q.Err != nil {
		q.Err = errors.Wrap(q.Err, "failed to get sql string for tx_id")
		return q
	}

	q.sql = q.sql.
		Where(fmt.Sprintf("hlc.tx_id in (%s)", sqlTxIDs), args...).
		// applying ordering to preserve natural order of ledger changes within op and ops within tx
		OrderBy("hlc.tx_id", "hlc.op_id", "hlc.order_number")

	return q
}
