package storage

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Matches - helper struct to store `matches`
type Matches struct {
	repo *pgdb.DB
}

// NewMatches - creates new instance of `Matches`
func NewMatches(repo *pgdb.DB) *Matches {
	return &Matches{
		repo: repo,
	}
}

func convertMatch(match history2.Match) []interface{} {
	return []interface{}{
		match.ID,
		match.OperationID,
		match.OfferID,
		match.BaseAmount,
		match.QuoteAmount,
		match.BaseAsset,
		match.QuoteAsset,
		match.Price,
		match.CreatedAt,
	}
}

// Insert - inserts new match
func (s *Matches) Insert(matches []history2.Match) error {
	columns := []string{
		"id",
		"operation_id",
		"offer_id",
		"base_amount",
		"quote_amount",
		"base_asset",
		"quote_asset",
		"price",
		"created_at",
	}

	err := matchesBatchInsert(s.repo, matches, "matches", columns, convertMatch)
	if err != nil {
		return errors.Wrap(err, "failed to insert matches")
	}

	return nil
}
