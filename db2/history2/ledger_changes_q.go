package history2

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// LedgerChangesQ is a helper struct to aid in configuring queries that loads
// ledger change structures.
type LedgerChangesQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewLedgerChangesQ - creates new instance of LedgerChangesQ
func NewLedgerChangesQ(repo *pgdb.DB) LedgerChangesQ {
	return LedgerChangesQ{
		repo: repo,
		selector: sq.Select(
			"ledger_changes.tx_id",
			"ledger_changes.op_id",
			"ledger_changes.order_number",
			"ledger_changes.effect",
			"ledger_changes.entry_type",
			"ledger_changes.payload",
		).From("ledger_changes"),
	}
}

// FilterByTransactionID - returns q with filter by transaction ID
func (q LedgerChangesQ) FilterByTransactionID(id int64) LedgerChangesQ {
	q.selector = q.selector.Where("ledger_changes.tx_id = ?", id)
	return q
}

// OrderByNumber - orders by order_number
func (q LedgerChangesQ) OrderByNumber() LedgerChangesQ {
	q.selector = q.selector.OrderBy("ledger_changes.order_number asc")
	return q
}

// Select - selects slice from the db, if no ledger changes found - returns nil, nil
func (q LedgerChangesQ) Select() ([]LedgerChanges, error) {
	var result []LedgerChanges
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger changes")
	}

	return result, nil
}
