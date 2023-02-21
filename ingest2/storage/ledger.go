package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Ledger - handles write operations on db level for ledgers
type Ledger struct {
	repo *pgdb.DB
}

// NewLedger - creates new instance of ledger
func NewLedger(repo *pgdb.DB) *Ledger {
	return &Ledger{
		repo: repo,
	}
}

// Insert - inserts Ledger into DB
func (s *Ledger) Insert(ledger *history2.Ledger) error {
	sql := squirrel.Insert("ledgers").Columns("id", "sequence", "hash", "previous_hash", "closed_at",
		"tx_count", "data").Values(ledger.ID, ledger.Sequence, ledger.Hash, ledger.PreviousHash, ledger.ClosedAt,
		ledger.TxCount, ledger.Data)

	err := s.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert ledger", logan.F{
			"ledger_seq": ledger.Sequence,
		})
	}

	return nil
}
