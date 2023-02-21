package core2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//FeesEmptyRole - defines is used as default in core when account role for which fee should be applied is not specified
// (when fee is set for account or global)
const FeesEmptyRole = 0

// FeesQ is a helper struct to aid in configuring queries that loads
// fee structs.
type FeesQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewFeesQ - creates new instance of Feesq
func NewFeesQ(repo *pgdb.DB) FeesQ {
	return FeesQ{
		repo: repo,
		selector: sq.Select("f.fee_type", "f.asset", "f.subtype", "f.fixed", "f.percent", "f.lastmodified",
			"f.account_id", "f.account_role", "f.lower_bound", "f.upper_bound", "f.hash").
			From("fee_state f"),
	}
}

// Page - returns Q with specified limit and offset params
func (q FeesQ) Page(params pgdb.OffsetPageParams) FeesQ {
	orderBys := []string{"f.hash " + params.Order, "f.lower_bound " + params.Order, "f.upper_bound " + params.Order}
	q.selector = params.ApplyTo(q.selector.OrderBy(orderBys...))
	return q
}

//FilterByAddress - returns q with filter by address
func (q FeesQ) FilterByAddress(address string) FeesQ {
	q.selector = q.selector.Where("f.account_id = ?", address)
	return q
}

//FilterByAsset - returns q with filter by asset
func (q FeesQ) FilterByAsset(asset string) FeesQ {
	q.selector = q.selector.Where("f.asset = ?", asset)
	return q
}

//FilterByType - returns q with filter by fee type
func (q FeesQ) FilterByType(feeType int32) FeesQ {
	q.selector = q.selector.Where("f.fee_type = ?", feeType)
	return q
}

//FilterBySubtype - returns q with filter by fee subtype
func (q FeesQ) FilterBySubtype(subtype int64) FeesQ {
	q.selector = q.selector.Where("f.subtype = ?", subtype)
	return q
}

//FilterByAccountRole - returns q with filter by account type
func (q FeesQ) FilterByAccountRole(role uint64) FeesQ {
	q.selector = q.selector.Where("f.account_role = ?", role)
	return q
}

//FilterGlobal - returns q with filter for global fees (where account_id and account_role are not set)
func (q FeesQ) FilterGlobal() FeesQ {
	q.selector = q.selector.Where("f.account_role =?", FeesEmptyRole).Where("f.account_id = ?", "")
	return q
}

//FilterByLowerBound - returns q with filter by lower bound
func (q FeesQ) FilterByLowerBound(lowerBound int64) FeesQ {
	q.selector = q.selector.Where("f.lower_bound = ?", lowerBound)
	return q
}

//FilterByUpperBound - returns q with filter by upper bound
func (q FeesQ) FilterByUpperBound(upperBound int64) FeesQ {
	q.selector = q.selector.Where("f.upper_bound = ?", upperBound)
	return q
}

//FilterByAmount - returns q with filter by upper bound
func (q FeesQ) FilterByAmount(amount int64) FeesQ {
	q.selector = q.selector.
		Where("f.upper_bound >= ?", amount).
		Where("f.lower_bound <= ?", amount)
	return q
}

// Get - loads a row from `fees`
// returns nil, nil - if fee does not exists
// returns error if more than one fee found
func (q FeesQ) Get() (*Fee, error) {
	var result Fee
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load fee")
	}

	return &result, nil
}

//Select - selects slice from the db, if no fees found - returns nil, nil
func (q FeesQ) Select() ([]Fee, error) {
	var result []Fee
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load fees")
	}

	return result, nil
}
