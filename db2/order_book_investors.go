package db2

import (
	"fmt"
	"strings"
)

type OrderBookInvestors struct {
	OrderBookId int64 `db:"order_book_id"`
	Quantity    int64 `db:"quantity"`
}

type OrderBooksInvestors []OrderBookInvestors

func (r OrderBookInvestors) String() string {
	return fmt.Sprintf("(%d,%d)", r.OrderBookId, r.Quantity)
}

func (r OrderBooksInvestors) String() string {
	var res string
	if len(r) == 0 {
		return "(0,0)"
	}

	for _, v := range r {
		res += fmt.Sprintf("%s,", v)
	}
	return strings.TrimSuffix(res, ",")
}
