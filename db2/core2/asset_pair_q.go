package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AssetPairsQ is a helper struct to aid in configuring queries that loads asset pairs
type AssetPairsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAssetPairsQ - creates new instance of AssetPairsQ with no filters
func NewAssetPairsQ(repo *pgdb.DB) AssetPairsQ {
	return AssetPairsQ{
		repo: repo,
		selector: sq.Select(
			"asset_pairs.base",
			"asset_pairs.quote",
			"asset_pairs.current_price",
			"asset_pairs.physical_price",
			"asset_pairs.physical_price_correction",
			"asset_pairs.max_price_step",
			"asset_pairs.policies",
		).From("asset_pair asset_pairs"),
	}
}

// SelectByAssets - selects slice of assets with FilterByAssets applied
func (q AssetPairsQ) SelectByAssets(bases, quotes []string) ([]AssetPair, error) {
	return q.FilterByAssets(bases, quotes).Select()
}

// FilterByBaseAsset - returns Q with filter by base
func (q AssetPairsQ) FilterByBaseAsset(base string) AssetPairsQ {
	q.selector = q.selector.Where("asset_pairs.base = ?", base)
	return q
}

// FilterByQuoteAsset - returns Q with filter by quote
func (q AssetPairsQ) FilterByQuoteAsset(quote string) AssetPairsQ {
	q.selector = q.selector.Where("asset_pairs.quote = ?", quote)
	return q
}

// FilterByAsset - returns Q with filter by asset (no matter base or quote)
func (q AssetPairsQ) FilterByAsset(code string) AssetPairsQ {
	q.selector = q.selector.Where("(asset_pairs.base = ? OR asset_pairs.quote = ?)", code, code)
	return q
}

// FilterByPolicy - returns Q with filter by policy
func (q AssetPairsQ) FilterByPolicy(mask uint64) AssetPairsQ {
	q.selector = q.selector.Where("asset_pairs.policies & ? = ?", mask, mask)
	return q
}

// FilterByAssets - filters pairs by baseAsset in baseAssets and quoteAsset in quoteAssets
func (q AssetPairsQ) FilterByAssets(baseAssets, quoteAssets []string) AssetPairsQ {
	q.selector = q.selector.Where(sq.Eq{"base": baseAssets}).Where(sq.Eq{"quote": quoteAssets})
	return q
}

// GetByBaseAndQuote - loads a row from `asset_pairs` found with matching base and quote assets
// returns nil, nil - if such pair doesn't exists
func (q AssetPairsQ) GetByBaseAndQuote(base, quote string) (*AssetPair, error) {
	return q.FilterByBaseAsset(base).FilterByQuoteAsset(quote).Get()
}

// Page - returns Q with specified limit and offset params
func (q AssetPairsQ) Page(params pgdb.OffsetPageParams) AssetPairsQ {
	q.selector = params.ApplyTo(q.selector, "asset_pairs.base", "asset_pairs.quote")
	return q
}

// WithBaseAsset - joins base asset
func (q AssetPairsQ) WithBaseAsset() AssetPairsQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "base_assets")...).
		LeftJoin("asset base_assets ON asset_pairs.base = base_assets.code")

	return q
}

// WithQuoteAsset - joins quote asset
func (q AssetPairsQ) WithQuoteAsset() AssetPairsQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "quote_assets")...).
		LeftJoin("asset quote_assets ON asset_pairs.quote = quote_assets.code")

	return q
}

// Get - loads a row from `asset_pairs`
// returns nil, nil - if asset pair does not exists
// returns error if more than one asset pair found
func (q AssetPairsQ) Get() (*AssetPair, error) {
	var result AssetPair
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load asset pair")
	}

	return &result, nil
}

// Select - selects slice from the db, if no pairs found - returns nil, nil
func (q AssetPairsQ) Select() ([]AssetPair, error) {
	var result []AssetPair
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load asset pairs")
	}

	return result, nil
}
