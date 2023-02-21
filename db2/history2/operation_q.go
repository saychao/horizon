package history2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var operationColumns = []string{"id", "tx_id", "type", "details",
	"ledger_close_time", "source"}

type OperationQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

func NewOperationQ(repo *pgdb.DB) OperationQ {
	return OperationQ{
		repo: repo,
		selector: sq.Select(
			"op.id",
			"op.tx_id",
			"op.type",
			"op.details",
			"op.ledger_close_time",
			"op.source").
			From("operations op"),
	}
}

func (q OperationQ) FilterByID(ids ...uint64) OperationQ {
	q.selector = q.selector.Where(sq.Eq{"op.id": ids})
	return q
}

func (q OperationQ) FilterByOperationsTypes(types []int) OperationQ {
	q.selector = q.selector.Where(sq.Eq{"op.type": types})
	return q
}

func (q OperationQ) FilterByOperationSource(source string) OperationQ {
	q.selector = q.selector.Where(sq.Eq{"op.source": source})
	return q
}

// Page - apply paging params to the query
func (q OperationQ) Page(pageParams pgdb.CursorPageParams) OperationQ {
	q.selector = pageParams.ApplyTo(q.selector, "op.id")
	return q
}

// Select - selects slice from the db, if no operations found - returns nil, nil
func (q OperationQ) Select() ([]Operation, error) {

	var result []Operation
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load operations")
	}

	return result, nil
}

// Get - loads a row
// returns nil, nil - if row does not exists
// returns error if more than one row found
func (q OperationQ) Get() (*Operation, error) {
	var result Operation
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load poll")
	}

	return &result, nil
}
