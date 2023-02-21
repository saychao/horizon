package core

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

var _ SaleAnteQI = &SaleAnteQ{}

type SaleAnteQI interface {
	// returns nil, nil if sale ante not found
	ByKey(balanceID string, saleID uint64) (*SaleAnte, error)
	// filters by sale id
	ForSale(saleID string) SaleAnteQI
	// filters by balance id
	ForBalance(balanceID string) SaleAnteQI
	// Select loads the results of the query specified by `q`
	Select() ([]SaleAnte, error)
}

type SaleAnteQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) SaleAntes() *SaleAnteQ {
	return &SaleAnteQ{
		parent: q,
		sql:    selectSaleAnte,
	}
}

func (q *SaleAnteQ) ByKey(balanceID string, saleID uint64) (*SaleAnte, error) {
	result := new(SaleAnte)
	query := selectSaleAnte.Limit(1).Where("sa.sale_id = ? AND sa.participant_balance_id = ?", saleID, balanceID)
	err := q.parent.Get(result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (q *SaleAnteQ) ForSale(saleID string) SaleAnteQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("sa.sale_id = ?", saleID)
	return q
}

func (q *SaleAnteQ) ForBalance(balanceID string) SaleAnteQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("sa.participant_balance_id = ?", balanceID)
	return q
}

func (q *SaleAnteQ) Select() ([]SaleAnte, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	result := make([]SaleAnte, 0)
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

var selectSaleAnte = sq.Select(
	"sa.sale_id",
	"sa.participant_balance_id",
	"sa.amount",
).From("sale_ante sa")
