// Package core contains database record definitions useable for
// reading rows from a Stellar Core db
package core

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/go/xdr"
)

// LedgerHeader is row of data from the `ledgerheaders` table
type LedgerHeader struct {
	LedgerHash     string           `db:"ledgerhash"`
	PrevHash       string           `db:"prevhash"`
	BucketListHash string           `db:"bucketlisthash"`
	CloseTime      int64            `db:"closetime"`
	Sequence       uint32           `db:"ledgerseq"`
	Version        uint64           `db:"version"`
	Data           xdr.LedgerHeader `db:"data"`
}

// Q is a helper struct on which to hang common queries against a stellar
// core database.
type Q struct {
	*pgdb.DB

	err error
	sql sq.SelectBuilder
}

func NewQ(repo *pgdb.DB) *Q {
	return &Q{
		DB: repo,
	}
}

func (q *Q) GetRepo() *pgdb.DB {
	return q.DB
}

func (q *Q) Exec(query sq.Sqlizer) error {
	return q.DB.Exec(query)
}

// Q interface helper for testing purposes mainly

type QInterface interface {
	GetRepo() *pgdb.DB
	// DEPRECATED
	LedgerHeaderBySequence(dest interface{}, seq int32) error
	// DEPRECATED
	ElderLedger(dest *int32) error
	// DEPRECATED
	LatestLedger(dest interface{}) error
	// DEPRECATED
	BalancesByAddress(dest interface{}, addy string) error
	// DEPRECATED
	BalanceByID(dest interface{}, bid string) error
	// DEPRECATED
	TransactionByHash(dest interface{}, hash string) error
	// DEPRECATED
	TransactionsByLedger(dest interface{}, seq int32) error
	// DEPRECATED
	TransactionFeesByLedger(dest interface{}, seq int32) error
	FeeEntries() FeeEntryQI
	Exec(query sq.Sqlizer) error // Was query (sqlx.Rows, error)
	// Returns nil, if not found
	FeeByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account, amount int64) (*FeeEntry, error)
	FeesByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account) ([]FeeEntry, error)

	//LimitsV2 - creates new limitsV2 query helper
	LimitsV2() LimitsV2QI
	//StatisticsV2 - creates new statisticsV2 query helper
	StatisticsV2() StatisticsV2QI
	// Accounts - creates new accounts query helper
	Accounts() AccountQI
	// Assets - creates new assets query helper
	Assets() AssetQI
	// AccountKyc - creates new account_kyc query helper
	AccountKYC() AccountKYCQI
	// Balances - creates new balances query builder
	Balances() BalancesQI
	//KeyValue - creates new KeyValue query helper
	KeyValue() KeyValueQI

	Trusts() *TrustQ
	Offers() *OfferQ
	OrderBook() *OrderBookQ
	Sales() *SaleQ
	SaleAntes() *SaleAnteQ

	// AssetPairs - creates new asset pair query helper
	AssetPairs() AssetPairsQ

	// ExternalSystemAccountID - returns builder to access external system account IDs
	ExternalSystemAccountID() ExternalSystemAccountIDQI
	ExternalSystemAccountIDPool() *ExternalSystemAccountIDPoolQ

	References() *ReferenceQ
}

// PriceLevel represents an aggregation of offers to trade at a certain
// price.
type PriceLevel struct {
	Pricen int32   `db:"pricen"`
	Priced int32   `db:"priced"`
	Pricef float64 `db:"pricef"`
	Amount int64   `db:"amount"`
}

// Transaction is row of data from the `txhistory` table from stellar-core
type Transaction struct {
	TransactionHash string                    `db:"txid"`
	LedgerSequence  int32                     `db:"ledgerseq"`
	Index           int32                     `db:"txindex"`
	Envelope        xdr.TransactionEnvelope   `db:"txbody"`
	Result          xdr.TransactionResultPair `db:"txresult"`
	ResultMeta      xdr.TransactionMeta       `db:"txmeta"`
}

// TransactionFee is row of data from the `txfeehistory` table from stellar-core
type TransactionFee struct {
	TransactionHash string                 `db:"txid"`
	LedgerSequence  int32                  `db:"ledgerseq"`
	Index           int32                  `db:"txindex"`
	Changes         xdr.LedgerEntryChanges `db:"txchanges"`
}

// ElderLedger represents the oldest "ingestable" ledger known to the
// stellar-core database this ingestion system is communicating with.  Horizon,
// which wants to operate on a contiguous range of ledger data (i.e. free from
// gaps) uses the elder ledger to start importing in the case of an empty
// database.
//
// Due to the design of stellar-core, ledger 1 will _always_ be in the database,
// even when configured to catchup minimally, so we cannot simply take
// MIN(ledgerseq). Instead, we can find whether or not 1 is the elder ledger by
// checking for the presence of ledger 2.
func (q *Q) ElderLedger(dest *int32) error {
	var found bool
	err := q.GetRaw(&found, `
		SELECT CASE WHEN EXISTS (
		    SELECT *
		    FROM ledgerheaders
		    WHERE ledgerseq = 2
		)
		THEN CAST(1 AS BIT)
		ELSE CAST(0 AS BIT) END
	`)

	if err != nil {
		return err
	}

	// if ledger 2 is present, use it 1 as the elder ledger (since 1 is guaranteed
	// to be present)
	if found {
		*dest = 1
		return nil
	}

	err = q.GetRaw(dest, `
		SELECT COALESCE(MIN(ledgerseq), 0)
		FROM ledgerheaders
		WHERE ledgerseq > 2
	`)

	return err
}

// LatestLedger loads the latest known ledger
func (q *Q) LatestLedger(dest interface{}) error {
	return q.GetRaw(dest, `SELECT COALESCE(MAX(ledgerseq), 0) FROM ledgerheaders`)
}

func (q *Q) Assets() AssetQI {
	return &assetQ{
		parent: q,
		sql:    selectAsset,
	}
}

func (q *Q) KeyValue() KeyValueQI {
	return &KeyValueQ{
		parent: q,
		sql:    selectKeyValue,
	}
}

func (q *Q) LimitsV2() LimitsV2QI {
	return &LimitsV2Q{
		parent: q,
		sql:    selectLimitsV2,
	}
}

func (q *Q) StatisticsV2() StatisticsV2QI {
	return &StatisticsV2Q{
		parent: q,
		sql:    selectStatisticsV2,
	}
}

// ExternalSystemAccountID - returns builder to access external system account IDs
func (q *Q) ExternalSystemAccountID() ExternalSystemAccountIDQI {
	return &externalSystemAccountIDQ{
		parent: q,
		sql:    selectExternalSystemAccountIDs,
	}
}

// AssetPairs - creates new asset pair query helper
func (q *Q) AssetPairs() AssetPairsQ {
	return &assetPairQ{
		parent: q,
		sql:    selectAssetPair,
	}
}

func (q *Q) ExternalSystemAccountIDPool() *ExternalSystemAccountIDPoolQ {
	return NewExternalSystemAccountIDPoolQ(q)
}
