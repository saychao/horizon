package history2

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// TransactionsQ is a helper struct to aid in configuring queries that loads
// transactions structures.
type TransactionsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewTransactionsQ - creates new instance of TransactionsQ
func NewTransactionsQ(repo *pgdb.DB) TransactionsQ {
	return TransactionsQ{
		repo: repo,
		selector: sq.Select(
			"transactions.id",
			"transactions.hash",
			"transactions.ledger_sequence",
			"transactions.ledger_close_time",
			"transactions.application_order",
			"transactions.account",
			"transactions.operation_count",
			"transactions.envelope",
			"transactions.result",
			"transactions.meta",
			"transactions.valid_after",
			"transactions.valid_before",
		).
			From("transactions").
			// To apply filters on ledger_changes properties:
			Distinct().
			LeftJoin("ledger_changes ON ledger_changes.tx_id = transactions.id"),
	}
}

// FilterByLedgerEntryTypes - returns q with filter by entry types
func (q TransactionsQ) FilterByLedgerEntryTypes(types ...int) TransactionsQ {
	q.selector = q.selector.Where(sq.Eq{"ledger_changes.entry_type": types})
	return q
}

// FilterByEffectTypes - returns q with filter by effect(ledger entry change) types
func (q TransactionsQ) FilterByEffectTypes(types ...int) TransactionsQ {
	q.selector = q.selector.Where(sq.Eq{"ledger_changes.effect": types})
	return q
}

// FilterByID - returns q with filter by transaction ID
func (q TransactionsQ) FilterByID(id uint64) TransactionsQ {
	q.selector = q.selector.Where("transactions.id = ?", id)
	return q
}

//FilterByHash - filters transaction by hash
func (q TransactionsQ) FilterByHash(hash string) TransactionsQ {
	q.selector = q.selector.Where("transactions.hash = ?", hash)
	return q
}

func (q TransactionsQ) FilterLedgerCloseTimeBefore(time time.Time) TransactionsQ {
	q.selector = q.selector.Where(sq.Lt{"transactions.ledger_close_time": time})
	return q
}

func (q TransactionsQ) FilterLedgerCloseTimeAfter(time time.Time) TransactionsQ {
	q.selector = q.selector.Where(sq.Gt{"transactions.ledger_close_time": time})
	return q
}

// GetByID loads a row from `transactions`, by ID
// returns nil, nil - if transaction does not exists
func (q TransactionsQ) GetByID(id uint64) (*Transaction, error) {
	return q.FilterByID(id).Get()
}

// Page - returns Q with specified limit and cursor params
func (q TransactionsQ) Page(params pgdb.CursorPageParams) TransactionsQ {
	q.selector = params.ApplyTo(q.selector, "transactions.id")
	return q
}

// Get - loads a row from `transactions`
// returns nil, nil - if transaction does not exists
// returns error if more than one Transaction found
func (q TransactionsQ) Get() (*Transaction, error) {
	var result Transaction
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load transaction")
	}

	return &result, nil
}

//Select - selects slice from the db, if no transactions found - returns nil, nil
func (q TransactionsQ) Select() ([]Transaction, error) {
	var result []Transaction

	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load transactions")
	}

	return result, nil
}

func (q TransactionsQ) GetByHash(hash string) (*Transaction, error) {
	q.selector = q.selector.Where("transactions.hash = ?", hash)
	return q.Get()
}
