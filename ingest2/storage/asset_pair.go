package storage

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

type AssetPair struct {
	repo *pgdb.DB
}

func NewAssetPair(repo *pgdb.DB) *AssetPair {
	return &AssetPair{
		repo: repo,
	}
}

func (q *AssetPair) InsertAssetPair(assetPair history2.AssetPair) error {
	sql := sq.Insert("asset_pairs").
		Columns("base", "quote", "current_price", "ledger_close_time").
		Values(assetPair.Base, assetPair.Quote, assetPair.CurrentPrice, assetPair.LedgerCloseTime)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert asset pair", logan.F{
			"base":  assetPair.Base,
			"quote": assetPair.Quote,
		})
	}

	return nil
}
