package history2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// MatchQ is a helper struct to aid in configuring queries that loads matches
type MatchQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewMatchQ returns new instance of MatchQ
func NewMatchQ(repo *pgdb.DB) MatchQ {
	return MatchQ{
		repo: repo,
		selector: sq.Select(
			"m.id",
			"m.operation_id",
			"m.offer_id",
			"m.base_amount",
			"m.quote_amount",
			"m.base_asset",
			"m.quote_asset",
			"m.price",
			"m.created_at",
		).From("matches m"),
	}
}

// FilterByAssetPair - returns Q with filter by asset pair
func (q MatchQ) FilterByAssetPair(base, quote string) MatchQ {
	q.selector = q.selector.
		Where("m.base_asset = ?", base).
		Where("m.quote_asset = ?", quote)
	return q
}

// Page - apply paging params to the query
func (q MatchQ) Page(pageParams pgdb.CursorPageParams) MatchQ {
	q.selector = pageParams.ApplyTo(q.selector, "m.id")
	return q
}

// Select - selects slice from the db, if no matches found - returns nil, nil
func (q MatchQ) Select() ([]Match, error) {

	var result []Match
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load matches")
	}

	return result, nil
}
