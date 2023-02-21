package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/saychao/horizon/db2"
)

// BalancesQ is a helper struct to aid in configuring queries that loads balances
type BalancesQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewBalancesQ - creates new instance of BalanceQ with no filters
func NewBalancesQ(repo *pgdb.DB) BalancesQ {
	return BalancesQ{
		repo: repo,
		selector: sq.Select("balances.balance_id", "balances.sequential_id", "balances.asset", "balances.account_id",
			"balances.amount", "balances.locked").From("balance balances"),
	}
}

// GetByAddress  - gets balance by address, if one does not exists - returns nil, nil
func (q BalancesQ) GetByAddress(address string) (*Balance, error) {
	return q.FilterByAddress(address).Get()
}

// FilterAssetOwner - returns new Q with added filter for asset owner
func (q BalancesQ) FilterByAssetOwner(ownerID string) BalancesQ {
	q.selector = q.selector.Join("asset a on a.code = balances.asset").Where("a.owner = ?", ownerID)
	return q
}

// FilterByAddress - returns new Q with added filter for balance address
func (q BalancesQ) FilterByAddress(address string) BalancesQ {
	q.selector = q.selector.Where("balances.balance_id = ?", address)
	return q
}

// FilterByAccount - returns new Q with added filter for account address
func (q BalancesQ) FilterByAccount(accountAddress string) BalancesQ {
	q.selector = q.selector.Where("balances.account_id = ?", accountAddress)
	return q
}

// FilterByAsset - returns new Q with added filter for asset
func (q BalancesQ) FilterByAsset(asset string) BalancesQ {
	q.selector = q.selector.Where("balances.asset = ?", asset)
	return q
}

// FilterByAmountLwOrEq - returns new Q with added filter for amount lower or equal then provided
func (q BalancesQ) FilterByAmountLwOrEq(amount uint64) BalancesQ {
	q.selector = q.selector.Where("balances.amount <= ?", amount)
	return q
}

// FilterByAmount - returns new Q with added filter for amount greater then provided
func (q BalancesQ) FilterByAmountGt(amount uint64) BalancesQ {
	q.selector = q.selector.Where("balances.amount > ?", amount)
	return q
}

// Page - returns Q with specified limit and offset params
func (q BalancesQ) Page(params pgdb.OffsetPageParams) BalancesQ {
	q.selector = params.ApplyTo(q.selector, "balances.balance_id")
	return q
}

// WithAsset - joins asset
func (q BalancesQ) WithAsset() BalancesQ {
	q.selector = q.selector.Columns(db2.GetColumnsForJoin(assetColumns, "assets")...).
		LeftJoin("asset assets ON balances.asset = assets.code")
	return q
}

// Get - selects balance from db using specified filters. Returns nil, nil - if one does not exists
// Returns error if more than one exists
func (q BalancesQ) Get() (*Balance, error) {
	var result Balance
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get balance")
	}

	return &result, nil
}

// Select - selects balances from fb using specified filters. Returns nil, nil - if one does not exists
func (q BalancesQ) Select() ([]Balance, error) {
	var result []Balance
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select balances")
	}

	return result, nil
}
