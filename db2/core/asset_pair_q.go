package core

import (
	sql2 "database/sql"

	sq "github.com/Masterminds/squirrel"
)

// AssetPairsQ is a helper interface to aid in configuring queries that loads
// slices or entry of Asset Pair structs.
type AssetPairsQ interface {
	// ByCode returns nil, if asset pair not found
	ByCode(base, quote string) (*AssetPair, error)
	// ForAssets loads pairs where baseAsset in baseAssets and quoteAsset in quoteAssets
	ForAssets(baseAssets, quoteAssets []string) AssetPairsQ
	// Select - selects all assets pairs with specified filters
	Select() ([]AssetPair, error)
}

// assetQ is a helper struct to aid in configuring queries that loads
// slices or entry of Asset structs.
type assetPairQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *assetPairQ) Select() ([]AssetPair, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var assetPairs []AssetPair
	err := q.parent.Select(&assetPairs, q.sql)
	if err == sql2.ErrNoRows {
		return nil, nil
	}

	return assetPairs, err
}

// returns nil, if not found
func (q *assetPairQ) ByCode(base, quote string) (*AssetPair, error) {
	sql := selectAssetPair.Where("base = ? AND quote = ?", base, quote)
	var result AssetPair
	err := q.parent.Get(&result, sql)
	if err == sql2.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

// ForAssets loads pairs where baseAsset in baseAssets and quoteAsset in quoteAssets
func (q *assetPairQ) ForAssets(baseAssets, quoteAssets []string) AssetPairsQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"base": baseAssets}).Where(sq.Eq{"quote": quoteAssets})
	return q
}

var selectAssetPair = sq.Select("a.base, a.quote, a.current_price, a.physical_price, a.physical_price_correction, a.max_price_step, a.policies").From("asset_pair a")
