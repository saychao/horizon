package core

import (
	sql2 "database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2/sqx"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AssetQI is a helper interface to aid in configuring queries that loads
// slices or entry of Asset structs.
type AssetQI interface {
	// ByCode - returns asset by code, if not found returns nil, nil
	ByCode(assetCode string) (*Asset, error)
	// ForOwner - filters assets by owner
	ForOwner(owner string) AssetQI
	// ForCodes - filters assets by code
	ForCodes(codes []string) AssetQI
	// ForPolicy -returns assets with specified policy
	ForPolicy(policy uint32) AssetQI
	// Select - selects assets for specified filter
	Select() ([]Asset, error)
}

// assetQ is a helper struct to aid in configuring queries that loads
// slices or entry of Asset structs.
type assetQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// ByCode - returns asset by code, if not found returns nil, nil
func (q *assetQ) ByCode(code string) (*Asset, error) {
	sql := selectAsset.Where("code = ?", code)
	var result Asset
	err := q.parent.Get(&result, sql)
	if err == sql2.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to load asset", map[string]interface{}{
			"code": code,
		})
	}
	return &result, err
}

// ForOwner - filters assets by owner
func (q *assetQ) ForOwner(ownerID string) AssetQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("a.owner = ?", ownerID)
	return q
}

func (q *assetQ) ForPolicy(policy uint32) AssetQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("a.policies & ? = ?", policy, policy)
	return q
}

// ForCodes - filters assets by code
func (q *assetQ) ForCodes(codes []string) AssetQI {
	if q.Err != nil {
		return q
	}

	query, values := sqx.InForString("code", codes...)
	q.sql = q.sql.Where(query, values...)
	return q
}

// Select - selects assets for specified filter
func (q *assetQ) Select() ([]Asset, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "failed before select. builder had error")
	}

	var result []Asset
	err := q.parent.Select(&result, q.sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select assets")
	}

	return result, nil
}

var selectAsset = sq.Select("a.code",
	"a.owner",
	"a.preissued_asset_signer",
	"a.details",
	"a.max_issuance_amount",
	"a.available_for_issueance",
	"a.issued",
	"a.policies",
	"a.pending_issuance",
	"a.trailing_digits",
	"a.type",
).From("asset a")
