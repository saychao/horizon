package history2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var assetColumns = []string{
	"code",
	"owner",
	"preissued_asset_signer",
	"details",
	"max_issuance_amount",
	"available_for_issuance",
	"issued",
	"pending_issuance",
	"policies",
	"trailing_digits",
	"type",
	"state",
}

// AssetQ is a helper struct to aid in configuring queries that loads assets
type AssetQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAssetQ- creates new instance of AssetQ
func NewAssetQ(repo *pgdb.DB) AssetQ {
	return AssetQ{
		repo: repo,
		selector: sq.Select(
			"asset.code",
			"asset.owner",
			"asset.details",
			"asset.preissued_asset_signer",
			"asset.max_issuance_amount",
			"asset.available_for_issuance",
			"asset.issued",
			"asset.pending_issuance",
			"asset.policies",
			"asset.trailing_digits",
			"asset.type",
			"asset.state",
		).From("asset asset"),
	}
}

// GetByCode - get asset by code
func (q AssetQ) GetByCode(code string) (*Asset, error) {
	q.selector = q.selector.Where("asset.code = ?", code)
	return q.Get()
}

//GetByOwner - gets asset by owner address, returns nil, nil if one does not exist
func (q AssetQ) ByOwner(address string) (*Asset, error) {
	q.selector = q.selector.Where("asset.owner = ?", address)
	return q.Get()
}

//Get - selects asset from db, returns nil, nil if one does not exists
func (q AssetQ) Get() (*Asset, error) {
	var result Asset
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load asset")
	}

	return &result, nil
}
