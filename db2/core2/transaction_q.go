package core2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// TransactionQ - helper struct to get transactions from db
type TransactionQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewTransactionQ - creates new instance of TransactionQ
func NewTransactionQ(repo *pgdb.DB) TransactionQ {
	return TransactionQ{
		repo: repo,
		selector: sq.Select(
			"tx.txid",
			"tx.ledgerseq",
			"tx.txindex",
			"tx.txbody",
			"tx.txresult",
			"tx.txmeta").
			From("txhistory tx"),
	}
}

// FilterByLedgerSeq returns TransactionQ filtered by ledger sequence
func (q TransactionQ) FilterByLedgerSeq(seq int32) TransactionQ {
	q.selector = q.selector.Where("ledgerseq = ?", seq)
	return q
}

// FilterByLedgerSeqRange returns TransactionQ filtered by range of ledger sequences, including boundaries.
func (q TransactionQ) FilterByLedgerSeqRange(fromSeq int32, toSeq int32) TransactionQ {
	q.selector = q.selector.Where("ledgerseq >= ? AND ledgerseq <= ?", fromSeq, toSeq)
	return q
}

// OrderByLedgerSeq - returns TransactionQ with ordered by sequence
func (q TransactionQ) OrderByLedgerSeq() TransactionQ {
	q.selector = q.selector.OrderBy("ledgerseq ASC")
	return q
}

// GetByLedgerRange returns ordered slice of transaction filtered by range of ledger sequences, including boundaries.
// Returns nil, nil if transactions do not exist
func (q TransactionQ) GetByLedgerRange(fromSeq int32, toSeq int32) ([]Transaction, error) {
	return q.FilterByLedgerSeqRange(fromSeq, toSeq).
		OrderByLedgerSeq().
		Select()
}

func (q TransactionQ) Get() (*Transaction, error) {
	var result Transaction
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load transaction")
	}
	return &result, nil
}

// Select - selects slice from the db, if no pairs found - returns nil, nil
func (q TransactionQ) Select() ([]Transaction, error) {
	var result []Transaction
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load transactions")
	}

	return result, nil
}

// GetByHash returns transaction for given hash. Returns nil, nil if there is no transaction with provided hash
func (q TransactionQ) GetByHash(hash string) (*Transaction, error) {
	q.selector = q.selector.Where("tx.txid = ?", hash)
	return q.Get()
}
