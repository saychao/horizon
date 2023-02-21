package core

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2"
)

type OrderBookQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) OrderBook() *OrderBookQ {
	return &OrderBookQ{
		parent: q,
		sql: sq.Select(
			"o.base_amount",
			"o.quote_amount",
			"o.price",
			"o.created_at",
		).From("offer o"),
	}
}

func (q *OrderBookQ) ForAssets(base, quote string) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("base_asset_code = ? AND quote_asset_code = ?", base, quote)
	return q
}

// ForOrderBookID - filters offers by order book id
func (q *OrderBookQ) ForOrderBookID(orderBookID uint64) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("order_book_id = ?", orderBookID)
	return q
}

func (q *OrderBookQ) Direction(isBuy bool) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	orderDirection := "ASC"
	if isBuy {
		orderDirection = "DESC"
	}

	q.sql = q.sql.Where("is_buy = ?", isBuy).OrderBy("price " + orderDirection).OrderBy("created_at ASC")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OrderBookQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

// InvestorsCount returns quantity of the unique investors for each order book.
func (q *OrderBookQ) InvestorsCount() (db2.OrderBooksInvestors, error) {
	dest := make([]db2.OrderBookInvestors, 0)

	query := sq.Select("o.order_book_id, count(*) as quantity").
		From("(SELECT DISTINCT owner_id, order_book_id from offer) AS o").
		GroupBy("o.order_book_id").
		OrderBy("quantity DESC")

	err := q.parent.Select(&dest, query)
	return dest, err
}
