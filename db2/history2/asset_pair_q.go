package history2

import (
	"database/sql"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
)

type AssetPairQ struct {
	repo *pgdb.DB
}

func NewAssetPairQ(repo *pgdb.DB) *AssetPairQ {
	return &AssetPairQ{
		repo: repo,
	}
}

func (q *AssetPairQ) AssetPairPriceAt(base, quote string, ts time.Time) (int64, error) {
	var result int64
	err := q.repo.Get(&result, sq.Select("ap.current_price").
		From("asset_pairs ap").
		Where("ap.base = ? AND ap.quote = ?", base, quote).
		Where("ap.ledger_close_time <= ?", ts).
		OrderBy("ledger_close_time DESC").
		Limit(1))
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, errors.Wrap(err, "failed to load asset pair", logan.F{"base": base, "quote": quote})
	}
	return result, nil
}
