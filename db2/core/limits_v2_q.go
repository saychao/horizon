package core

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

var selectLimitsV2 = sq.Select("lim.id",
	"lim.account_type",
	"lim.account_id",
	"lim.stats_op_type",
	"lim.asset_code",
	"lim.is_convert_needed",
	"lim.daily_out",
	"lim.weekly_out",
	"lim.monthly_out",
	"lim.annual_out").From("limits_v2 lim")

type LimitsV2QI interface {
	// Global - filters to select only global limits
	Global() LimitsV2QI
	// ForAccountID - filters to select only for specific accountID
	ForAccountID(accountID string) LimitsV2QI
	// ForAccountType - filters to select only for account type
	ForAccountType(accountType int32) LimitsV2QI
	// ForAsset - filters to select only for specified asset
	ForAsset(asset string) LimitsV2QI
	// ForAsset - filters to select only for specified stats type
	ForStatsOpType(statsType int32) LimitsV2QI
	// ForAccountByAccountType - selects limit for account
	ForAccount(account *Account) ([]LimitsV2Entry, error)
	// Select - selects using build query
	Select() ([]LimitsV2Entry, error)
}

type LimitsV2Q struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// ForAccount - selects limit for account
func (q *LimitsV2Q) ForAccount(account *Account) ([]LimitsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	// to make sure that there is no sql injection
	var accID xdr.AccountId
	err := accID.SetAddress(account.AccountID)
	if err != nil {
		return nil, errors.From(errors.New("Invalid account ID"), map[string]interface{}{
			"acc_id": account.AccountID,
		})
	}
	accountIDStr := fmt.Sprintf("'%s'", account.AccountID)

	query := fmt.Sprintf("select distinct on (stats_op_type, asset_code, is_convert_needed)  id, "+
		"account_type, account_id, stats_op_type, asset_code, is_convert_needed, daily_out, "+
		"weekly_out, monthly_out, annual_out "+
		"from limits_v2 "+
		"where (account_type=%d or account_type is null) and (account_id=%s or account_id is null)"+
		"order by stats_op_type, asset_code, is_convert_needed, account_id = %s, "+
		"account_type = %d asc", account.RoleID, accountIDStr, accountIDStr, account.RoleID)

	var result []LimitsV2Entry
	err = q.parent.SelectRaw(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 for account")
	}

	return result, nil
}

// Global - filters to select only global limits
func (q *LimitsV2Q) Global() LimitsV2QI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("lim.account_id is NULL AND lim.stats_op_type is NULL")
	return q
}

// ForAccountID - filters to select only for specific accountID
func (q *LimitsV2Q) ForAccountID(accountID string) LimitsV2QI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("lim.account_id = ?", accountID)
	return q
}

// ForAccountType - filters to select only for account type
func (q *LimitsV2Q) ForAccountType(accountType int32) LimitsV2QI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("lim.account_type = ?", accountType)
	return q
}

// ForAsset - filters to select only for specified asset
func (q *LimitsV2Q) ForAsset(asset string) LimitsV2QI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("lim.asset_code = ?", asset)
	return q
}

// ForAsset - filters to select only for specified stats type
func (q *LimitsV2Q) ForStatsOpType(statsType int32) LimitsV2QI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("lim.stats_op_type = ?", statsType)
	return q
}

// Select - selects using build query
func (q *LimitsV2Q) Select() ([]LimitsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []LimitsV2Entry
	err := q.parent.Select(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load limits_v2 by stats op type")
	}
	return result, nil
}
