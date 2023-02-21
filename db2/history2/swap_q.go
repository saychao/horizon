package history2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// SwapsQ is a helper struct to aid in configuring queries that loads
// swap structures.
type SwapsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewSwapsQ - creates new instance of SwapsQ
func NewSwapsQ(repo *pgdb.DB) SwapsQ {
	return SwapsQ{
		repo: repo,
		selector: sq.Select(
			"s.id",
			"s.created_at",
			"s.lock_time",
			"s.source_account",
			"s.destination_account",
			"s.source_balance",
			"s.destination_balance",
			"s.secret",
			"s.secret_hash",
			"s.amount",
			"s.asset",
			"s.source_percent_fee",
			"s.source_fixed_fee",
			"s.destination_fixed_fee",
			"s.destination_percent_fee",
			"s.details",
			"s.state",
		).From("swaps s"),
	}
}

// GetByID loads a row from `swaps`, by ID
// returns nil, nil - if swap does not exists
func (q SwapsQ) GetByID(id int64) (*Swap, error) {
	q.selector = q.selector.Where("s.id = ?", id)
	return q.Get()
}

// FilterBySwapState - returns q with filter by state
func (q SwapsQ) FilterByState(state int32) SwapsQ {
	q.selector = q.selector.Where("s.state = ?", state)
	return q
}

// FilterByAsset - returns q with filter by asset
func (q SwapsQ) FilterByAsset(asset string) SwapsQ {
	q.selector = q.selector.Where("s.asset = ?", asset)
	return q
}

// FilterBySourceBalance - returns q with filter by source balance
func (q SwapsQ) FilterBySourceBalance(sourceBalance string) SwapsQ {
	q.selector = q.selector.Where("s.source_balance = ?", sourceBalance)
	return q
}

// FilterBySource - returns q with filter by source account
func (q SwapsQ) FilterBySource(source string) SwapsQ {
	q.selector = q.selector.Where("s.source_account = ?", source)
	return q
}

// FilterByDestinationBalance - returns q with filter by destination balance
func (q SwapsQ) FilterByDestinationBalance(destinationBalance string) SwapsQ {
	q.selector = q.selector.Where("s.destination_balance = ?", destinationBalance)
	return q
}

// FilterByDestination - returns q with filter by destination account
func (q SwapsQ) FilterByDestination(destination string) SwapsQ {
	q.selector = q.selector.Where("s.destination_account = ?", destination)
	return q
}

// Page - returns Q with specified limit and offset params
func (q SwapsQ) Page(params pgdb.CursorPageParams) SwapsQ {
	q.selector = params.ApplyTo(q.selector, "s.id")
	return q
}

// Get - loads a row from `swaps`
// returns nil, nil - if swap does not exists
// returns error if more than one CreateSwap found
func (q SwapsQ) Get() (*Swap, error) {
	var result Swap
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load swap")
	}

	return &result, nil
}

//Select - selects slice from the db, if no swaps found - returns nil, nil
func (q SwapsQ) Select() ([]Swap, error) {
	var result []Swap
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load swaps")
	}

	return result, nil
}
