package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var assetColumns = []string{"code", "owner", "preissued_asset_signer", "details",
	"max_issuance_amount", "available_for_issueance", "issued",
	"pending_issuance", "policies", "trailing_digits", "type", "state"}

//AssetsQ - helper struct to load assets from db
type AssetsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAssetsQ - returns new instance of AssetsQ
func NewAssetsQ(repo *pgdb.DB) AssetsQ {
	return AssetsQ{
		repo:     repo,
		selector: sq.Select(assetColumns...).From("asset assets"),
	}
}

// GetByCode - loads a row from `assets` found with matching code
// returns nil, nil - if such asset doesn't exists
func (q AssetsQ) GetByCode(code string) (*Asset, error) {
	return q.FilterByCode(code).Get()
}

//SelectByPolicy - selects slice of assets by policy mask
func (q AssetsQ) SelectByPolicy(mask uint64) ([]Asset, error) {
	return q.FilterByPolicy(mask).Select()
}

// FilterByCode - returns q with filter by code
func (q AssetsQ) FilterByCode(code string) AssetsQ {
	q.selector = q.selector.Where("assets.code = ?", code)
	return q
}

// FilterByCodes - returns q with filter by codes
func (q AssetsQ) FilterByCodes(codes []string) AssetsQ {
	q.selector = q.selector.Where(sq.Eq{"assets.code": codes})
	return q
}

// FilterByTypes - returns q with filter by types
func (q AssetsQ) FilterByTypes(types []uint64) AssetsQ {
	q.selector = q.selector.Where(sq.Eq{"assets.type": types})
	return q
}

// FilterByOwner - returns q with filter by owner ID
func (q AssetsQ) FilterByOwner(ownerID string) AssetsQ {
	q.selector = q.selector.Where("assets.owner = ?", ownerID)
	return q
}

// FilterByState - returns q with filter by owner ID
func (q AssetsQ) FilterByState(state uint32) AssetsQ {
	q.selector = q.selector.Where("assets.state = ?", state)
	return q
}

// FilterByPolicy - returns q with filter by policy
func (q AssetsQ) FilterByPolicy(mask uint64) AssetsQ {
	q.selector = q.selector.Where("assets.policies & ? = ?", mask, mask)
	return q
}

// Page - returns Q with specified limit and offset params
func (q AssetsQ) Page(params pgdb.OffsetPageParams) AssetsQ {
	q.selector = params.ApplyTo(q.selector, "assets.code")
	return q
}

// Select - selects slice from the db, if no assets found - returns nil, nil
func (q AssetsQ) Select() ([]Asset, error) {
	var result []Asset
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load assets")
	}

	return result, nil
}

// Get - loads a row from `assets`
// returns nil, nil - if asset does not exists
// returns error if more than one asset found
func (q AssetsQ) Get() (*Asset, error) {
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
