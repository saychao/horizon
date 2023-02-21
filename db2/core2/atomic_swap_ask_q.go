package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AtomicSwapAskQ is a helper struct to aid in configuring queries that loads atomic swap bids
type AtomicSwapAskQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewAtomicSwapAskQ - creates new instance of AtomicSwapAskQ with no filters
func NewAtomicSwapAskQ(repo *pgdb.DB) AtomicSwapAskQ {
	return AtomicSwapAskQ{
		repo: repo,
		selector: sq.Select(
			"b.id",
			"b.owner_id",
			"b.base_asset_code",
			"b.base_balance_id",
			"b.base_amount",
			"b.locked_amount",
			"b.is_cancelled",
			"b.details",
			"b.created_at",
		).From("atomic_swap_ask b"),
	}
}

// GetByID - loads a row from `atomic_swap_bids` found with matching id
// returns nil, nil - if such atomic swap bid doesn't exists
func (q AtomicSwapAskQ) GetByID(id uint64) (*AtomicSwapAsk, error) {
	return q.FilterByID(id).Get()
}

// FilterByCode - returns q with filter by code
func (q AtomicSwapAskQ) FilterByID(id uint64) AtomicSwapAskQ {
	q.selector = q.selector.Where("b.id = ?", id)
	return q
}

// FilterByIDs - returns q with filter by ids
func (q AtomicSwapAskQ) FilterByIDs(ids []uint64) AtomicSwapAskQ {
	q.selector = q.selector.Where(sq.Eq{"b.id": ids})
	return q
}

// FilterByOwner - returns q with filter by owner ID
func (q AtomicSwapAskQ) FilterByOwner(ownerID string) AtomicSwapAskQ {
	q.selector = q.selector.Where("b.owner_id = ?", ownerID)
	return q
}

func (q AtomicSwapAskQ) IDSelector() AtomicSwapAskQ {
	q.selector = sq.Select("b.id").From("atomic_swap_ask b")
	return q
}

func (q AtomicSwapAskQ) FilterByBaseAssets(codes []string) AtomicSwapAskQ {
	q.selector = q.selector.Where(sq.Eq{"b.base_asset_code": codes})
	return q
}

func (q AtomicSwapAskQ) FilterByQuoteAssets(codes []string) AtomicSwapAskQ {
	q.selector = q.selector.LeftJoin("atomic_swap_quote_asset qa ON qa.ask_id = b.id").
		Where(sq.Eq{"qa.quote_asset": codes})
	return q
}

// Page - returns Q with specified limit and offset params
func (q AtomicSwapAskQ) Page(params pgdb.OffsetPageParams) AtomicSwapAskQ {
	q.selector = params.ApplyTo(q.selector, "b.id")
	return q
}

// Get - loads a row from `atomic_swap_bids`
// returns nil, nil - if atomic swap bid does not exists
// returns error if more than one asset found
func (q AtomicSwapAskQ) Get() (*AtomicSwapAsk, error) {
	var result AtomicSwapAsk
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load atomic swap ask")
	}

	return &result, nil
}

func (q *AtomicSwapAskQ) Select() ([]AtomicSwapAsk, error) {
	var result []AtomicSwapAsk
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load atomic swap asks")
	}

	return result, nil
}

func (q AtomicSwapAskQ) SelectIDs() ([]uint64, error) {
	var result []uint64
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load atomic swap ask ids")
	}

	return result, nil
}
