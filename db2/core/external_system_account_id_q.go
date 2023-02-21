package core

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ExternalSystemAccountIDQI - provides methods to select from db
type ExternalSystemAccountIDQI interface {
	// ForAccount - filters EXS accounts by accountID
	ForAccount(accountID string) ExternalSystemAccountIDQI
	// Select - selects slice of EXS account IDs
	Select() ([]ExternalSystemAccountID, error)
}

type externalSystemAccountIDQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// ForAccount - filters EXS accounts by accountID
func (q *externalSystemAccountIDQ) ForAccount(accountID string) ExternalSystemAccountIDQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("esa.account_id = ?", accountID)
	return q
}

// Select - selects slice of EXS account IDs
func (q *externalSystemAccountIDQ) Select() ([]ExternalSystemAccountID, error) {
	if q.Err != nil {
		err := errors.Wrap(q.Err, "failed to select due to error in query builder")
		return nil, err
	}

	var result []ExternalSystemAccountID
	err := q.parent.Select(&result, q.sql)
	if err != nil {
		return nil, errors.Wrap(err, "select failed")
	}

	return result, err
}

var selectExternalSystemAccountIDs = sq.Select(
	"esa.account_id",
	"esa.external_system_type",
	"esa.data",
	"p.expires_at as pool_entry_expires_at").
	From("external_system_account_id esa").
	LeftJoin("external_system_account_id_pool p on (p.account_id = esa.account_id AND " +
		"p.external_system_type = esa.external_system_type)")
