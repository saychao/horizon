package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AccountRoleQ is a helper struct to aid in configuring queries that loads
// accountRole structs.
type SignerRuleQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAccountRoleQ - creates new instance of AccountRoleQ
func NewSignerRuleQ(repo *pgdb.DB) SignerRuleQ {
	return SignerRuleQ{
		repo: repo,
		selector: sq.Select("sr.id",
			"sr.resource",
			"sr.action",
			"sr.forbids",
			"sr.is_default",
			"sr.owner_id",
			"sr.details",
		).From("signer_rules sr"),
	}
}

//FilterByRole - filter rules by role ID
func (q SignerRuleQ) FilterByRole(roleID uint64) SignerRuleQ {
	q.selector = q.selector.Join("signer_role_rules srr on srr.rule_id = sr.id").Where("srr.role_id = ?", roleID)
	return q
}

//FilterByID - returns q with filter by address
func (q SignerRuleQ) FilterByID(id uint64) SignerRuleQ {
	q.selector = q.selector.Where("sr.id = ?", id)
	return q
}

// Page - returns Q with specified limit and offset params
func (q SignerRuleQ) Page(params pgdb.OffsetPageParams) SignerRuleQ {
	q.selector = params.ApplyTo(q.selector, "sr.id")
	return q
}

// Get - loads a row from `account_roles`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRole found
func (q SignerRuleQ) Get() (*SignerRule, error) {
	var result SignerRule
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load signer rule")
	}

	return &result, nil
}

//Select - selects slice from the db, if no account rules found - returns nil, nil
func (q SignerRuleQ) Select() ([]SignerRule, error) {
	var result []SignerRule
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load signer rules")
	}

	return result, nil
}
