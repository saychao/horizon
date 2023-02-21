package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type LimitsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewLimitsQ - default constructor for LimitsQ which
// creates LimitsQ with given pgdb.DB and default selector
func NewLimitsQ(repo *pgdb.DB) LimitsQ {
	return LimitsQ{
		repo: repo,
		selector: sq.
			Select("limits.id, " +
				"limits.account_type, " +
				"limits.account_id, " +
				"limits.stats_op_type, " +
				"limits.asset_code," +
				"limits.is_convert_needed, " +
				"limits.daily_out, " +
				"limits.weekly_out, " +
				"limits.monthly_out, " +
				"limits.annual_out").
			From("limits_v2 limits"),
	}
}

// FilterByAccount - adds accountID filter for query to Limits table
func (q LimitsQ) FilterByAccount(accountID string) LimitsQ {
	q.selector = q.selector.Where("limits.account_id = ?", accountID)
	return q
}

//FilterByAccountRole - returns q with filter by account role
func (q LimitsQ) FilterByAccountRole(role uint64) LimitsQ {
	q.selector = q.selector.Where("limits.account_type = ?", role)
	return q
}

//General
func (q LimitsQ) General() LimitsQ {
	q.selector = q.selector.Where("limits.account_type is null and limits.account_id is null")
	return q
}

// Page - returns Q with specified limit and offset params
func (q LimitsQ) Page(params pgdb.OffsetPageParams) LimitsQ {
	q.selector = params.ApplyTo(q.selector, "limits.id")
	return q
}

//FilterByAsset - returns q with filter by asset
func (q LimitsQ) FilterByAsset(asset string) LimitsQ {
	q.selector = q.selector.Where("limits.asset_code = ?", asset)
	return q
}

//FilterByStatsOpType - returns q with filter by stats op type
func (q LimitsQ) FilterByStatsOpType(statsOpType int32) LimitsQ {
	q.selector = q.selector.Where("limits.stats_op_type = ?", statsOpType)
	return q
}

// Select - loads rows from `limits_v2`
// returns nil, nil - if limits for particular account does not exists
func (q LimitsQ) Select() ([]Limits, error) {
	var result []Limits
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select limits")
	}

	return result, nil
}
