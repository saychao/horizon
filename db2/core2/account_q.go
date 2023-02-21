package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AccountsQ is a helper struct to aid in configuring queries that loads
// account structs.
type AccountsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAccountsQ - creates new instance of AccountsQ
func NewAccountsQ(repo *pgdb.DB) AccountsQ {
	return AccountsQ{
		repo: repo,
		selector: sq.Select("accounts.account_id",
			"accounts.sequential_id",
			"accounts.referrer",
			"accounts.sequential_id",
			"accounts.role_id",
		).From("accounts accounts"),
	}
}

// GetByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q AccountsQ) GetByAddress(address string) (*Account, error) {
	return q.FilterByAddress(address).Get()
}

//FilterByAddress - returns q with filter by address
func (q AccountsQ) FilterByAddress(address string) AccountsQ {
	q.selector = q.selector.Where("accounts.account_id = ?", address)
	return q
}

//FilterByAddresses - returns q with filter by addresses
func (q AccountsQ) FilterByAddresses(addresses ...string) AccountsQ {
	q.selector = q.selector.Where(sq.Eq{"accounts.account_id": addresses})
	return q
}

//FilterByReferrer - returns q with filter by referrer
func (q AccountsQ) FilterByReferrer(address string) AccountsQ {
	q.selector = q.selector.Where("accounts.referrer = ?", address)
	return q
}

// FilterByRole - returns q with filter by role
func (q AccountsQ) FilterByRole(ids ...uint64) AccountsQ {
	q.selector = q.selector.Where(sq.Eq{"accounts.role_id": ids})
	return q
}

// Page - returns Q with specified limit and offset params
func (q AccountsQ) Page(params pgdb.OffsetPageParams) AccountsQ {
	q.selector = params.ApplyTo(q.selector, "accounts.sequential_id")
	return q
}

// Get - loads a row from `accounts`
// returns nil, nil - if account does not exists
// returns error if more than one Account found
func (q AccountsQ) Get() (*Account, error) {
	var result Account
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}

//Select - selects slice from the db, if no accounts found - returns nil, nil
func (q AccountsQ) Select() ([]Account, error) {
	var result []Account
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account accounts")
	}

	return result, nil
}
