package history

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type PricePoint struct {
	Price     float64   `db:"price" json:"price"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}

var selectPricePoint = sq.Select("hp.price", "hp.timestamp").From("history_price hp")

func (q *Q) PriceHistory(base, quote string, since time.Time) ([]PricePoint, error) {
	var result []PricePoint

	stmt := `
		select sum(price)/count(price) as price, max(timestamp) as timestamp
		from (
			select price, timestamp, ntile(360) over (order by "timestamp") as bucket
			from history_price
			where timestamp > $1 and base_asset = $2 and quote_asset = $3
		) as t
		group by bucket order by bucket`

	err := q.DB.SelectRaw(&result, stmt, since, base, quote)
	return result, err
}

func (q *Q) LastPrice(base, quote string) (*PricePoint, error) {
	var result PricePoint

	sqq := selectPricePoint.
		Where("base_asset = ?", base).
		Where("quote_asset = ?", quote).
		Limit(1).
		OrderBy("timestamp desc")

	err := q.Get(&result, sqq)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}
