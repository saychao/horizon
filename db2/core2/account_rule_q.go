package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AccountRuleQ is a helper struct to aid in configuring queries that loads
// accountRule structs.
type AccountRuleQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAccountRuleQ - creates new instance of AccountRuleQ
func NewAccountRuleQ(repo *pgdb.DB) AccountRuleQ {
	return AccountRuleQ{
		repo: repo,
		selector: sq.Select("ar.id",
			"ar.resource",
			"ar.action",
			"ar.forbids",
			"ar.details",
		).From("account_rules ar"),
	}
}

//FilterByRole - filter rules by role ID
func (q AccountRuleQ) FilterByRole(roleID uint64) AccountRuleQ {
	q.selector = q.selector.Join("account_role_rules arr on arr.rule_id = ar.id").Where("arr.role_id = ?", roleID)
	return q
}

//FilterByID - filters account rules by id
func (q AccountRuleQ) FilterByID(id uint64) AccountRuleQ {
	q.selector = q.selector.Where("ar.id = ?", id)
	return q
}

// Page - returns Q with specified limit and offset params
func (q AccountRuleQ) Page(params pgdb.OffsetPageParams) AccountRuleQ {
	q.selector = params.ApplyTo(q.selector, "ar.id")
	return q
}

// Get - loads a row from `account_rules`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRule found
func (q AccountRuleQ) Get() (*AccountRule, error) {
	var result AccountRule
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account rule")
	}

	return &result, nil
}

//Select - selects slice from the db, if no account rules found - returns nil, nil
func (q AccountRuleQ) Select() ([]AccountRule, error) {
	var result []AccountRule
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account rules")
	}

	return result, nil
}
