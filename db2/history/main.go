// Package history contains database record definitions useable for
// reading rows from a the history portion of horizon's database
package history

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

// EffectType is the numeric type for an effect, used as the `type` field in the
// `history_effects` table.
type EffectType int

// Q is a helper struct on which to hang common queries against a history
// portion of the horizon database.
type Q struct {
	*pgdb.DB
}

func (q *Q) GetRepo() *pgdb.DB {
	return q.DB
}

type QInterface interface {
	GetRepo() *pgdb.DB

	ElderLedger(dest interface{}) error
	LatestLedger(dest interface{}) error
	OldestOutdatedLedgers(dest interface{}, currentVersion int) error

	// Accounts
	Accounts() AccountsQI
	AccountByAddress(dest interface{}, addy string) error
	AccountByID(dest interface{}, id int64) error

	Balances() BalancesQI
	BalanceByID(dest interface{}, id string) error

	// Ledgers
	Ledgers() LedgersQI
	LedgerBySequence(dest interface{}, seq int32) error

	// Operations
	Operations() OperationsQI
	OperationByID(dest interface{}, id int64) error

	// Transactions
	Transactions() TransactionsQI
	TransactionByHash(dest interface{}, hash string) error
	TransactionByHashOrID(dest interface{}, hash string) error

	// PendingTransactions

	// prices history
	PriceHistory(base, quote string, since time.Time) ([]PricePoint, error)
	LastPrice(base, quote string) (*PricePoint, error)

	Trades() TradesQI

	// Sales - returns query builder for sales
	Sales() SalesQ

	// ReviewableRequests - provides builder of request to access reviewable requests
	ReviewableRequests() ReviewableRequestQI

	// LedgerChanges - provides builder to access ledger changes
	LedgerChanges() LedgerChangesQI

	//Contracts
	Contracts() ContractQI
	ContractsDetails() ContractsDetailsQI
	ContractDispute() ContractDisputeQI

	// Offers - provides builder to work with offer entries
	Offers() OffersQI

	// GetOldestLedgerSeq - returns oldest ledger sequence
	GetOldestLedgerSeq() (int32, error)
	// GetLatestLedgerSeq - returns latest ledger sequence available in DB
	GetLatestLedgerSeq() (int32, error)
}

// ReviewableRequests - provides builder of request to access reviewable requests
func (q *Q) ReviewableRequests() ReviewableRequestQI {
	return &ReviewableRequestQ{
		parent: q,
		sql:    selectReviewableRequest,
	}
}

// ElderLedger loads the oldest ledger known to the history database
func (q *Q) ElderLedger(dest interface{}) error {
	return q.GetRaw(dest, `SELECT COALESCE(MIN(sequence), 0) FROM history_ledgers`)
}

// LatestLedger loads the latest known ledger
func (q *Q) LatestLedger(dest interface{}) error {
	return q.GetRaw(dest, `SELECT COALESCE(MAX(sequence), 0) FROM history_ledgers`)
}

// OldestOutdatedLedgers populates a slice of ints with the first million
// outdated ledgers, based upon the provided `currentVersion` number
func (q *Q) OldestOutdatedLedgers(dest interface{}, currentVersion int) error {
	return q.SelectRaw(dest, `
		SELECT sequence
		FROM history_ledgers
		WHERE importer_version < $1
		ORDER BY sequence ASC
		LIMIT 1000000`, currentVersion)
}

// GetLatestLedgerSeq - returns latest ledger sequence available in DB
func (q *Q) GetLatestLedgerSeq() (int32, error) {
	var result int32
	err := q.LatestLedger(&result)
	return result, err
}

// GetOldestLedgerSeq - returns oldest ledger sequence
func (q *Q) GetOldestLedgerSeq() (int32, error) {
	var result int32
	err := q.ElderLedger(&result)
	return result, err
}
