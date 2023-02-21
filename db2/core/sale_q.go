package core

import (
	sq "github.com/Masterminds/squirrel"
)

type SaleQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Sales() *SaleQ {
	return &SaleQ{
		parent: q,
		sql:    selectSale,
	}
}

func (q *SaleQ) Select() ([]Sale, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	result := make([]Sale, 0)
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

var selectSale = sq.Select(
	"s.id",
).From("sale s")
