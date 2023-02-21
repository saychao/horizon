package history

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var selectOffers = sq.Select(
	"id",
	"offer_id",
	"owner_id",
	"base_asset",
	"quote_asset",
	"is_buy",
	"initial_base_amount",
	"current_base_amount",
	"price",
	"is_canceled",
	"created_at",
).From("history_offer")

type OffersQI interface {
	// ForBase filters offers by base asset
	ForBase(base string) OffersQI
	// ForQuote filters offers by quote asset
	ForQuote(quote string) OffersQI
	// ForOwnerID filters offers by owner id
	ForOwnerID(ownerID string) OffersQI
	// ForIsBuy filters offers by is buy
	ForIsBuy(isBuy bool) OffersQI
	// NoMatches filters offers that have not been matched with others
	NoMatches() OffersQI
	// PartiallyMatched filters offers that have been matched not fully with others
	PartiallyMatched() OffersQI
	// FullyMatched filters offers that have been matched fully with others
	FullyMatched() OffersQI
	// Canceled filters offers that have been canceled by owner
	Canceled(isCanceled bool) OffersQI
	// Page - applies page params
	Page(page db2.PageQuery) OffersQI
	// Select - selects offers by specifics filters
	Select() ([]Offer, error)
	// Insert - add new offer to database
	Insert(offer Offer) error
	// UpdateBaseAmount changes current base amount by offer id
	UpdateBaseAmount(baseAmount, offerID int64) error
	// Cancel set true to is_canceled field by offer id
	Cancel(offerID int64) error
}

type OffersQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Offers() OffersQI {
	return &OffersQ{
		parent: q,
		sql:    selectOffers,
	}
}

func (q *OffersQ) ForBase(base string) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"base_asset": base})
	return q
}

func (q *OffersQ) ForQuote(quote string) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"quote_asset": quote})
	return q
}

func (q *OffersQ) ForOwnerID(ownerID string) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"owner_id": ownerID})
	return q
}

func (q *OffersQ) ForIsBuy(isBuy bool) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"is_buy": isBuy})
	return q
}

func (q *OffersQ) NoMatches() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("initial_base_amount = current_base_amount")
	return q
}

func (q *OffersQ) PartiallyMatched() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("initial_base_amount != current_base_amount").
		Where("current_base_amount > 0")
	return q
}

func (q *OffersQ) FullyMatched() OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("current_base_amount = 0")
	return q
}

func (q *OffersQ) Canceled(isCanceled bool) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"is_canceled": isCanceled})
	return q
}

func (q *OffersQ) Page(page db2.PageQuery) OffersQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "id")
	return q
}

func (q *OffersQ) Select() ([]Offer, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "failed to build query before select")
	}

	var offers []Offer
	err := q.parent.Select(&offers, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to execute select query")
	}

	return offers, nil
}

func (q *OffersQ) Insert(offer Offer) error {
	query := sq.Insert("history_offer").SetMap(sq.Eq{
		"offer_id":            offer.OfferID,
		"owner_id":            offer.OwnerID,
		"base_asset":          offer.BaseAsset,
		"quote_asset":         offer.QuoteAsset,
		"is_buy":              offer.IsBuy,
		"initial_base_amount": offer.InitialBaseAmount,
		"current_base_amount": offer.CurrentBaseAmount,
		"price":               offer.Price,
		"is_canceled":         offer.IsCanceled,
		"created_at":          offer.CreatedAt,
	})

	err := q.parent.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query")
	}

	return nil
}

func (q *OffersQ) UpdateBaseAmount(baseAmount, offerID int64) error {
	query := sq.Update("history_offer").SetMap(sq.Eq{
		"current_base_amount": baseAmount,
	}).Where(sq.Eq{"offer_id": offerID})

	err := q.parent.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}

func (q *OffersQ) Cancel(offerID int64) error {
	query := sq.Update("history_offer").
		Set("is_canceled", true).
		Where(sq.Eq{"offer_id": offerID})

	err := q.parent.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}
