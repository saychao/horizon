package history

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2"
)

type TradesQI interface {
	ForPair(base, quote string) TradesQI
	Page(page db2.PageQuery) TradesQI
	// ForOrderBook - filters trades by order book
	ForOrderBook(orderBookID uint64) TradesQI
	Select(dest interface{}) error
}

type TradesQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Trades() TradesQI {
	return &TradesQ{
		parent: q,
		sql:    selectTrades,
	}
}

// ForOrderBook - filters trades by order book
func (q *TradesQ) ForOrderBook(orderBookID uint64) TradesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("order_book_id = ?", orderBookID)
	return q
}

func (q *TradesQ) ForPair(base, quote string) TradesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ht.base_asset = ? AND ht.quote_asset = ?", base, quote)
	return q
}

func (q *TradesQ) Page(page db2.PageQuery) TradesQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ht.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *TradesQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

var selectTrades = sq.Select(
	"ht.id",
	"ht.base_asset",
	"ht.quote_asset",
	"ht.base_amount",
	"ht.quote_amount",
	"ht.price",
	"ht.created_at").
	From("history_trades ht")
