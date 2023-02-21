package core2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AtomicSwapQuoteAssetQ is a helper to aid in configuring queries
// that loads slices or entry of AtomicSwapBidQuoteAsset structs.
type AtomicSwapQuoteAssetQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

func NewAtomicSwapQuoteAssetQ(repo *pgdb.DB) AtomicSwapQuoteAssetQ {
	return AtomicSwapQuoteAssetQ{
		repo: repo,
		selector: sq.Select(
			"qa.ask_id",
			"qa.quote_asset",
			"qa.price",
		).From("atomic_swap_quote_asset qa"),
	}
}

func (q AtomicSwapQuoteAssetQ) FilterByIDs(bidIDs []int64) AtomicSwapQuoteAssetQ {
	q.selector = q.selector.Where(sq.Eq{"qa.ask_id": bidIDs})
	return q
}

func (q AtomicSwapQuoteAssetQ) FilterByCodes(quoteAsset []string) AtomicSwapQuoteAssetQ {
	q.selector = q.selector.Where(sq.Eq{"qa.quote_asset": quoteAsset})
	return q
}

func (q AtomicSwapQuoteAssetQ) Select() ([]AtomicSwapQuoteAsset, error) {
	var result []AtomicSwapQuoteAsset
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select atomic swap quote assets")
	}

	return result, nil
}

func (q AtomicSwapQuoteAssetQ) Get() (*AtomicSwapQuoteAsset, error) {
	var result AtomicSwapQuoteAsset
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load atomic swap ask quote asset")
	}

	return &result, nil
}
