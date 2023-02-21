package history2

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/saychao/horizon/db2"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

// SalesQ is a helper struct to aid in configuring queries that loads
// sales structures.
type SalesQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewSalesQ - creates new instance of SalesQ
func NewSalesQ(repo *pgdb.DB) SalesQ {
	return SalesQ{
		repo: repo,
		selector: sq.Select(
			"sales.id",
			"sales.soft_cap",
			"sales.hard_cap",
			"sales.base_current_cap",
			"sales.base_hard_cap",
			"sales.sale_type",
			"sales.owner_address",
			"sales.base_asset",
			"sales.default_quote_asset",
			"sales.start_time",
			"sales.end_time",
			"sales.details",
			"sales.state",
			"sales.quote_assets",
			"sales.version",
			"sales.access_definition_type",
		).From("sales sales"),
	}
}

// Whitelisted - returns q with filter by sales whitelisted for address
func (q SalesQ) Whitelisted(address string) SalesQ {
	//selects sale ids for which account was explicitly whitelisted
	explicitlyAllowed, _, _ := sq.StatementBuilder.
		Select("(sr.key#>>'{sale,saleID}')::int").
		Distinct().From("account_specific_rules sr").
		Where("sr.address = ?").
		Where("not sr.forbids").ToSql()

	//selects sale ids with global whitelist
	globalAllowed, _, _ := sq.StatementBuilder.
		Select("(sr.key#>>'{sale,saleID}')::int").
		Distinct().From("account_specific_rules sr").
		Where("sr.address is null").
		Where("not sr.forbids").ToSql()

	//selects sale ids for which account was explicitly blacklisted
	explicitlyForbidden, _, _ := sq.StatementBuilder.
		Select("(sr.key#>>'{sale,saleID}')::int").
		Distinct().From("account_specific_rules sr").
		Where("sr.address = ?", address).
		Where("sr.forbids").ToSql()

	/*
		query go get whitelisted sale ids
		whitelisted = explicitly allowed + (globally allowed - explicitly forbidden)
		will give us exactly sales with account whitelisted
	*/
	whitelist := fmt.Sprintf(
		`(sales.id in ( %s union (%s except %s)))`,
		explicitlyAllowed,
		globalAllowed,
		explicitlyForbidden)

	q.selector = q.selector.
		Where(
			sq.Or{
				//we have 2 placeholder symbols in `whitelist` query (see above), thus we need 2 arguments (both address)
				sq.Expr(whitelist, address, address),
				sq.Expr("sales.version < ?", int32(xdr.LedgerVersionAddSaleWhitelists)),
			}).
		//Exclude sale owner
		Where("sales.owner_address != ?", address)
	return q
}

// FilterByID - returns q with filter by sale ID
func (q SalesQ) FilterByID(id uint64) SalesQ {
	q.selector = q.selector.Where("sales.id = ?", id)
	return q
}

// FilterByIDs - returns q with filter by ids
func (q SalesQ) FilterByIDs(ids []uint64) SalesQ {
	q.selector = q.selector.Where(sq.Eq{"sales.id": ids})
	return q
}

func (q SalesQ) WithAsset() SalesQ {
	q.selector = q.selector.Columns(db2.GetColumnsForJoin(assetColumns, "asset")...).
		LeftJoin("asset asset ON asset.code = sales.base_asset")
	return q
}

// GetByID loads a row from `sales`, by ID
// returns nil, nil - if sale does not exists
func (q SalesQ) GetByID(id uint64) (*Sale, error) {
	return q.FilterByID(id).Get()
}

// FilterByOwner - returns q with filter by Owner
func (q SalesQ) FilterByOwner(ownerAddress string) SalesQ {
	q.selector = q.selector.Where("sales.owner_address = ?", ownerAddress)
	return q
}

// FilterByBaseAsset - returns q with filter by BaseAsset
func (q SalesQ) FilterByBaseAsset(baseAssetCode string) SalesQ {
	q.selector = q.selector.Where("sales.base_asset = ?", baseAssetCode)
	return q
}

// FilterByAssetType - returns q with filter by AssetType
func (q SalesQ) FilterByAssetType(assetType uint64) SalesQ {
	q.selector = q.selector.Columns(db2.GetColumnsForJoin(assetColumns, "asset")...).
		LeftJoin("asset asset ON asset.code = sales.base_asset").
		Where("asset.type = ?", assetType)
	return q
}

// FilterByMaxEndTime - returns q with filter by max end time
func (q SalesQ) FilterByMaxEndTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.end_time <= ?", time)
	return q
}

// FilterByMaxStartTime - returns q with filter by start_time
func (q SalesQ) FilterByMaxStartTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.start_time <= ?", time)
	return q
}

// FilterByMinStartTime - returns q with filter by start_time
func (q SalesQ) FilterByMinStartTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.start_time >= ?", time)
	return q
}

// FilterByMinEndTime - returns q with filter by end_time
func (q SalesQ) FilterByMinEndTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.end_time >= ?", time)
	return q
}

// FilterByState - returns q with filter by state
func (q SalesQ) FilterByState(state uint64) SalesQ {
	q.selector = q.selector.Where("sales.state = ?", state)
	return q
}

// FilterBySaleType - returns q with filter by type
func (q SalesQ) FilterBySaleType(saleType uint64) SalesQ {
	q.selector = q.selector.Where("sales.sale_type = ?", saleType)
	return q
}

// FilterByMinHardCap - returns q with filter by min hard cap
func (q SalesQ) FilterByMinHardCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.hard_cap >= ?", value)
	return q
}

// FilterByMinSoftCap - returns q with filter by min soft cap
func (q SalesQ) FilterByMinSoftCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.soft_cap >= ?", value)
	return q
}

// FilterByMaxHardCap - returns q with filter by max hard cap
func (q SalesQ) FilterByMaxHardCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.hard_cap <= ?", value)
	return q
}

func (q SalesQ) FilterByParticipant(participant string, saleIDs []int64) SalesQ {
	q.selector = q.selector.LeftJoin("(SELECT DISTINCT (effect#>>'{matched,order_book_id}')::int sale_id, asset_code, account_id FROM participant_effects) pe on (sales.id = pe.sale_id and sales.state = ? and pe.asset_code = sales.base_asset)", regources.SaleStateClosed).
		LeftJoin("accounts a ON pe.account_id = a.id").
		Where(sq.Or{sq.Eq{"a.address": participant}, sq.Eq{"sales.id": saleIDs}}).Where("sales.owner_address != ?", participant)
	return q
}

// FilterByMaxSoftCap - returns q with filter by max sof cap
func (q SalesQ) FilterByMaxSoftCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.soft_cap <= ?", value)
	return q
}

// FilterByCompanyIDs - returns q with filter by company id. Must use this method with WithAsset only
func (q SalesQ) FilterByCompanyIDs(ids []uint64) SalesQ {
	q.selector = q.selector.
		Where(sq.Eq{"asset.details->>'company_id'": ids})
	return q
}

// Page - returns Q with specified limit and offset params
func (q SalesQ) Page(params pgdb.OffsetPageParams) SalesQ {
	q.selector = params.ApplyTo(q.selector, "sales.id")
	return q
}

// CursorPage - returns Q with specified limit and offset params
func (q SalesQ) CursorPage(params pgdb.CursorPageParams) SalesQ {
	q.selector = params.ApplyTo(q.selector, "sales.id")
	return q
}

// Get - loads a row from `sales`
// returns nil, nil - if sale does not exists
// returns error if more than one Sale found
func (q SalesQ) Get() (*Sale, error) {
	var result Sale
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load sale")
	}

	return &result, nil
}

// Select - selects slice from the db, if no sales found - returns nil, nil
func (q SalesQ) Select() ([]Sale, error) {
	var result []Sale
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sales")
	}

	return result, nil
}
