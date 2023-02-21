package storage

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

// ParticipantEffect - helper struct to store `operation participants`
type AccountSpecificRules struct {
	repo *pgdb.DB
}

// NewAccountSpecificRules - creates new instance of `AccountSpecificRules`
func NewAccountSpecificRules(repo *pgdb.DB) *AccountSpecificRules {
	return &AccountSpecificRules{
		repo: repo,
	}
}

//Insert - stores account specific rule into db
func (q *AccountSpecificRules) Insert(rule history2.AccountSpecificRule) error {
	columns := []string{"id", "address", "entry_type", "forbids", "key"}
	sql := sq.Insert("account_specific_rules").
		Columns(columns...).
		Values(rule.ID, rule.Address, rule.EntryType, rule.Forbids, rule.Key)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert account specific rule", logan.F{"rule_id": rule.ID})
	}

	return nil
}

//Remove - removes account specific rule from db
func (q *AccountSpecificRules) Remove(ruleID uint64) error {
	sql := sq.Delete("account_specific_rules").Where("id = ?", ruleID)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to remove account specific rule", logan.F{"rule_id": ruleID})
	}

	return nil
}
