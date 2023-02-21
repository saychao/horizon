package history2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// BalancesQ is a helper struct to aid in configuring queries that loads balances
type BalancesQ struct {
	repo *pgdb.DB
}

// NewBalancesQ - creates new instance of BalanceQ
func NewBalancesQ(repo *pgdb.DB) BalancesQ {
	return BalancesQ{
		repo: repo,
	}
}

// GetByAddress loads a row from `balances`, by address
// returns nil, nil - if balance does not exists
func (q BalancesQ) GetByAddress(address string) (*Balance, error) {
	var result Balance
	err := q.repo.Get(&result, sq.Select(balanceColumns...).
		From("balances balances").Where("balances.address = ?", address))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load balance by address", logan.F{"balance_address": address})
	}

	return &result, nil
}

// SelectByAddress loads rows from 'balances' by address
// returns nil, nil - if balances for given account does not exist
func (q BalancesQ) SelectByAddress(address ...string) ([]Balance, error) {
	var result []Balance

	err := q.repo.Select(&result, sq.Select(balanceColumns...).
		From("balances balances").
		Where(sq.Eq{
			"balances.address": address,
		}),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load balance by balance address", logan.F{"balances": address})
	}

	return result, nil
}

var balanceColumns = []string{"balances.id", "balances.account_id", "balances.address", "balances.asset_code"}
