package history2

import (
	"database/sql"

	"github.com/SafeRE-IT/horizon/db2/core2"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// SaleParticipationQ is a helper struct to aid in configuring queries that load
// sale participation structures from `participant_effects` table`.
type SaleParticipationQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewSaleParticipationQ - creates new instance of SaleParticipationQ
func NewSaleParticipationQ(repo *pgdb.DB) SaleParticipationQ {
	return SaleParticipationQ{
		repo: repo,
		selector: sq.Select().
			From("participant_effects pe").
			Join("accounts a ON pe.account_id = a.id").
			Where("(pe.effect#>>'{type}')::int = ?", EffectTypeMatched),
	}
}

// FilterByQuoteAsset - returns q with filter by quote asset
func (q SaleParticipationQ) FilterByQuoteAsset(asset string) SaleParticipationQ {
	q.selector = q.selector.Where("pe.effect#>>'{matched,charged,asset_code}' = ?", asset)
	return q
}

// FilterByParticipant - returns q with filter by participant
func (q SaleParticipationQ) FilterByParticipant(id string) SaleParticipationQ {
	q.selector = q.selector.Where("a.address = ?", id)
	return q
}

// FilterBySaleParams - returns q with filter by sale params
func (q SaleParticipationQ) FilterBySaleParams(id uint64, baseAsset, owner string) SaleParticipationQ {
	q.selector = q.selector.
		Where("(pe.effect#>>'{matched,order_book_id}')::int = ?", id).
		Where("pe.asset_code = ?", baseAsset).
		Where("a.address != ?", owner)

	return q
}

// FilterBySaleBaseAssets - returns q with filter by sales default quote assets
func (q SaleParticipationQ) FilterBySaleBaseAssets(baseAssets ...string) SaleParticipationQ {
	q.selector = q.selector.Where(sq.Eq{
		"pe.asset_code": baseAssets,
	})

	return q
}

// FilterByNotSaleOwners - returns q with filter by not sales owners
func (q SaleParticipationQ) FilterByNotSaleOwners(owners ...string) SaleParticipationQ {
	q.selector = q.selector.Where(sq.NotEq{
		"a.address": owners,
	})

	return q
}

// FilterBySaleIDs - returns q with filter by sale IDs
func (q SaleParticipationQ) FilterBySaleIDs(ids ...uint64) SaleParticipationQ {
	q.selector = q.selector.
		Where(sq.Eq{
			"(pe.effect#>>'{matched,order_book_id}')::int": ids,
		})

	return q
}

// Page - returns Q with specified cursor params
func (q SaleParticipationQ) Page(params pgdb.CursorPageParams) SaleParticipationQ {
	q.selector = params.ApplyTo(q.selector, "pe.id")
	return q
}

// Select - selects slice from the db, if no sales found - returns nil, nil
func (q SaleParticipationQ) Select() ([]SaleParticipation, error) {
	var result []SaleParticipation
	q.selector = q.selector.Columns("(pe.effect#>>'{matched,offer_id}')::int id",
		"a.address participant_id",
		"(pe.effect#>>'{matched,order_book_id}')::int sale_id",
		"pe.effect#>>'{matched,charged,asset_code}' quote_asset",
		"pe.effect#>>'{matched,charged,amount}' quote_amount",
		"pe.effect#>>'{matched,funded,asset_code}' base_asset",
		"pe.effect#>>'{matched,funded,amount}' base_amount")
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale participations")
	}

	return result, nil
}

// SelectParticipationsCount - returns slice of participations count with corresponding sale ID
func (q SaleParticipationQ) SelectParticipationsCount() ([]core2.SaleIDParticipantsCount, error) {
	var result []core2.SaleIDParticipantsCount

	q.selector = q.selector.Columns("(pe.effect#>>'{matched,order_book_id}')::int sale_id", "COUNT(*) p_count").
		GroupBy("sale_id")

	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to load sale participations")
	}

	return result, nil
}

// SelectParticipantsCount - returns slice of participants count with corresponding sale ID
func (q SaleParticipationQ) SelectParticipantsCount() ([]core2.SaleIDParticipantsCount, error) {
	var result []core2.SaleIDParticipantsCount

	q.selector = q.selector.Columns("(pe.effect#>>'{matched,order_book_id}')::int sale_id", "COUNT(DISTINCT pe.account_id) p_count").
		GroupBy("sale_id")

	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to load sale participants")
	}

	return result, nil
}
