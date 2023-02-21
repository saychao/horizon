package history2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AccountSpecificRulesQ is a helper struct to aid in configuring queries that loads accounts
type AccountSpecificRulesQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAccountSpecificRulesQ - creates new instance of AccountSpecificRulesQ
func NewAccountSpecificRulesQ(repo *pgdb.DB) AccountSpecificRulesQ {
	return AccountSpecificRulesQ{
		repo: repo,
		selector: sq.Select(
			"sr.id",
			"sr.address",
			"sr.entry_type",
			"sr.forbids",
			"sr.key",
		).From("account_specific_rules sr"),
	}
}

// FilterByAddress returns q with filter by address
func (q AccountSpecificRulesQ) FilterByAddress(address string) AccountSpecificRulesQ {
	q.selector = q.selector.Where("sr.address = ?", address)
	return q
}

// FilterByPermission returns q with filter by `forbids`
func (q AccountSpecificRulesQ) FilterByPermission(forbids bool) AccountSpecificRulesQ {
	q.selector = q.selector.Where("sr.forbids = ?", forbids)
	return q
}

// FilterBySale - returns q with filter by sale id
func (q AccountSpecificRulesQ) FilterBySale(saleID uint64) AccountSpecificRulesQ {
	q.selector = q.selector.Where("sr.key#>>'{sale,saleID}' = ?", saleID)
	return q
}

//Get - selects account from db, returns nil, nil if one does not exists
func (q AccountSpecificRulesQ) Get() (*AccountSpecificRule, error) {
	var result AccountSpecificRule
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account specific rule")
	}

	return &result, nil
}

// Page - returns Q with specified limit and offset params
func (q AccountSpecificRulesQ) Page(params pgdb.CursorPageParams) AccountSpecificRulesQ {
	q.selector = params.ApplyTo(q.selector, "sr.id")
	return q
}

//Get - selects account from db, returns nil, nil if one does not exists
func (q AccountSpecificRulesQ) Select() ([]AccountSpecificRule, error) {
	var result []AccountSpecificRule
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account specific rules")
	}

	return result, nil
}
