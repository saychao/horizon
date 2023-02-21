package core

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type ReferenceQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) References() *ReferenceQ {
	return &ReferenceQ{
		parent: q,
		sql:    selectReference,
	}
}

func (q *ReferenceQ) ForAccount(accountID string) *ReferenceQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("sender = ?", accountID)
	return q
}

// ByReference matches references by substring
func (q *ReferenceQ) ByReference(reference string) *ReferenceQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("reference ilike ?", fmt.Sprintf("%%%s%%", reference))
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *ReferenceQ) Select() ([]Reference, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	result := make([]Reference, 0)
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

var selectReference = sq.Select(
	"r.reference",
).From("reference r")
