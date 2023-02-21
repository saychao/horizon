package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ExternalSystemIDsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewExternalSystemIDsQ - default constructor for ExternalSystemIDsQ which
// creates ExternalSystemIDsQ with given pgdb.DB and default selector
func NewExternalSystemIDsQ(repo *pgdb.DB) ExternalSystemIDsQ {
	return ExternalSystemIDsQ{
		repo: repo,
		selector: sq.Select(
			"ext_pool.id",
			"ext_pool.external_system_type",
			"ext_pool.data",
			"ext_pool.account_id",
			"ext_pool.is_deleted",
			"to_timestamp(ext_pool.expires_at) at time zone 'UTC' expires_at",
			"to_timestamp(ext_pool.binded_at) at time zone 'UTC' binded_at").
			From("external_system_account_id_pool ext_pool"),
	}
}

// FilterByAccount - adds accountID filter for query to external system IDs table
func (esid ExternalSystemIDsQ) FilterByAccount(accountID string) ExternalSystemIDsQ {
	esid.selector = esid.selector.Where("ext_pool.account_id = ?", accountID)
	return esid
}

// Select - loads a rows from `external_system_account_id` left joining
// with table external_system_account_id_pool on account_id
// returns nil, nil - if external system IDs for particular account does not exists
func (esid ExternalSystemIDsQ) Select() ([]ExternalSystemID, error) {
	var result []ExternalSystemID
	err := esid.repo.Select(&result, esid.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select external system IDs")
	}

	return result, nil
}
