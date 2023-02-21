package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AccountRoleQ is a helper struct to aid in configuring queries that loads
// accountRole structs.
type AccountRoleQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAccountRoleQ - creates new instance of AccountRoleQ
func NewAccountRoleQ(repo *pgdb.DB) AccountRoleQ {
	return AccountRoleQ{
		repo: repo,
		selector: sq.Select("ar.id",
			"ar.details",
		).From("account_roles ar"),
	}
}

//FilterByID - returns q with filter by id of account role
func (q AccountRoleQ) FilterByID(id uint64) AccountRoleQ {
	q.selector = q.selector.Where("ar.id = ?", id)
	return q
}

// Page - returns Q with specified limit and offset params
func (q AccountRoleQ) Page(params pgdb.OffsetPageParams) AccountRoleQ {
	q.selector = params.ApplyTo(q.selector, "ar.id")
	return q
}

// Get - loads a row from `account_roles`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRole found
func (q AccountRoleQ) Get() (*AccountRole, error) {
	var result AccountRole
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account role")
	}

	return &result, nil
}

//Select - selects slice from the db, if no account roles found - returns nil, nil
func (q AccountRoleQ) Select() ([]AccountRole, error) {
	var result []AccountRole
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account account roles")
	}

	return result, nil
}
