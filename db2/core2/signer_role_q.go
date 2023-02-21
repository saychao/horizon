package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AccountRoleQ is a helper struct to aid in configuring queries that loads
// accountRole structs.
type SignerRoleQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewSignerRoleQ - creates new instance of AccountRoleQ
func NewSignerRoleQ(repo *pgdb.DB) SignerRoleQ {
	return SignerRoleQ{
		repo: repo,
		selector: sq.Select("sr.id",
			"sr.owner_id",
			"sr.details",
		).From("signer_roles sr"),
	}
}

// GetByID loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q SignerRoleQ) GetByID(id uint64) (*SignerRole, error) {
	return q.FilterByID(id).Get()
}

//FilterByID - returns q with filter by address
func (q SignerRoleQ) FilterByID(id uint64) SignerRoleQ {
	q.selector = q.selector.Where("sr.id = ?", id)
	return q
}

// Page - returns Q with specified limit and offset params
func (q SignerRoleQ) Page(params pgdb.OffsetPageParams) SignerRoleQ {
	q.selector = params.ApplyTo(q.selector, "sr.id")
	return q
}

// Get - loads a row from `account_roles`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRole found
func (q SignerRoleQ) Get() (*SignerRole, error) {
	var result SignerRole
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load signer role")
	}

	return &result, nil
}

//Select - selects slice from the db, if no account roles found - returns nil, nil
func (q SignerRoleQ) Select() ([]SignerRole, error) {
	var result []SignerRole
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load signer roles")
	}

	return result, nil
}
