package core

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2"
)

const (
	secondaryMarketID = 0
)

type OfferQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Offers() *OfferQ {
	return &OfferQ{
		parent: q,
		sql:    selectOffer,
	}
}

func (q *OfferQ) ForOfferID(offerID string) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("offer_id = ?", offerID)
	return q
}

func (q *OfferQ) ForAccount(accountID string) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("owner_id = ?", accountID)
	return q
}

func (q *OfferQ) ForOrderBookID(orderBookID uint64) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("order_book_id = ?", orderBookID)
	return q
}

func (q *OfferQ) OnlyPrimaryMarket() *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("order_book_id <> ?", secondaryMarketID)
	return q
}

func (q *OfferQ) ForAssets(base, quote string) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("base_asset_code = ? AND quote_asset_code = ?", base, quote)
	return q
}

func (q *OfferQ) IsBuy(isBuy bool) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("is_buy = ?", isBuy)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *OfferQ) Page(page db2.PageQuery) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "o.offer_id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OfferQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

func (q *OfferQ) OrderByPrice(desc bool) *OfferQ {
	if q.Err != nil {
		return q
	}

	order := "ASC"
	if desc {
		order = "DESC"
	}

	q.sql = q.sql.OrderBy("o.price " + order)
	return q
}

var selectOffer = sq.Select(
	"o.owner_id",
	"o.offer_id",
	"o.base_asset_code",
	"o.quote_asset_code",
	"o.is_buy",
	"o.base_amount",
	"o.quote_amount",
	"o.price",
	"o.fee",
	"o.base_balance_id",
	"o.quote_balance_id",
	"o.created_at",
	"o.order_book_id",
).From("offer o")
