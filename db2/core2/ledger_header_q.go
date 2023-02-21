package core2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// LedgerHeaderQ - helper struct to get ledger headers from db
type LedgerHeaderQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewLedgerHeaderQ - creates new instance of LedgerHeaderQ
func NewLedgerHeaderQ(repo *pgdb.DB) *LedgerHeaderQ {
	return &LedgerHeaderQ{
		repo: repo,
		selector: sq.Select(
			"l.ledgerhash",
			"l.prevhash",
			"l.bucketlisthash",
			"l.closetime",
			"l.ledgerseq",
			"l.version",
			"l.data",
		).From("ledgerheaders l"),
	}
}

// SelectBySequence returns []core.LedgerHeader by its sequence. Returns nil, nil if ledgerHeader does not exists
func (q *LedgerHeaderQ) SelectBySequence(seq []int32) ([]LedgerHeader, error) {
	query := q.selector.Where(sq.Eq{"ledgerseq": seq})
	var header []LedgerHeader
	err := q.repo.Select(&header, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger by sequence", logan.F{"seq": seq})
	}

	return header, nil
}

// GetBySequence returns *core.LedgerHeader by its sequence. Returns nil, nil if ledgerHeader does not exists
func (q *LedgerHeaderQ) GetBySequence(seq int32) (*LedgerHeader, error) {
	query := q.selector.Where("ledgerseq = ?", seq)
	var header LedgerHeader
	err := q.repo.Get(&header, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load ledger by sequence", logan.F{"seq": seq})
	}

	return &header, nil
}

// GetBySequenceRange returns ordered slice of ledger headers inside specified range of sequences, including boundaries.
// Returns nil, nil if ledgerHeaders do not exist
func (q *LedgerHeaderQ) GetBySequenceRange(fromSeq int32, toSeq int32) ([]LedgerHeader, error) {
	query := q.selector.Where("ledgerseq >= ? AND ledgerseq <= ?", fromSeq, toSeq).
		OrderBy("ledgerseq ASC")
	var headers []LedgerHeader
	err := q.repo.Select(&headers, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load ledger by sequence", logan.F{"fromSeq": fromSeq, "toSeq": toSeq})
	}

	return headers, nil
}

// GetLatestLedgerSeq - returns latest ledger sequence available in DB
func (q *LedgerHeaderQ) GetLatestLedgerSeq() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT COALESCE(MAX(ledgerseq), 1) FROM ledgerheaders")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get latest ledger seq")
	}

	return result, nil
}

// GetOldestLedgerSeq - returns oldest ledger sequence (which does not have gaps) available in the core db.
// Due to design of core it will always have 1 ledger in the db. So we try to find min ledger sequence > 1
// In case it's 2 we will return 1.
// Note: this method does not handle and should not handle existence of 2 gaps
func (q *LedgerHeaderQ) GetOldestLedgerSeq() (int32, error) {
	var result int32
	err := q.repo.GetRaw(&result, "SELECT COALESCE(MIN(ledgerseq), 1) FROM ledgerheaders WHERE ledgerseq > 1")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get oldest ledger seq")
	}

	if result == 2 {
		return 1, nil
	}

	return result, nil
}
