package core

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type CoreQMock struct {
	mock.Mock
}

func (q *CoreQMock) BalanceByID(dest interface{}, bid string) error {
	args := q.Called(dest, bid)
	return args.Error(0)
}

func (q *CoreQMock) Balances() BalancesQI {
	args := q.Called()
	return args.Get(0).(BalancesQI)
}

func (q *CoreQMock) BalancesByAddress(dest interface{}, addy string) error {
	return nil
}

func (q *CoreQMock) FeesByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account) ([]FeeEntry, error) {
	args := q.Called(feeType, asset, subtype, account)
	return args.Get(0).([]FeeEntry), args.Error(1)
}

func (q *CoreQMock) ExternalSystemAccountID() ExternalSystemAccountIDQI {
	args := q.Called()
	return args.Get(0).(ExternalSystemAccountIDQI)
}

func (q *CoreQMock) Offers() *OfferQ {
	args := q.Called()
	return args.Get(0).(*OfferQ)
}

func (q *CoreQMock) OrderBook() *OrderBookQ {
	args := q.Called()
	return args.Get(0).(*OrderBookQ)
}
func (q *CoreQMock) Trusts() *TrustQ {
	args := q.Called()
	return args.Get(0).(*TrustQ)
}

func (q *CoreQMock) FeeByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account, amount int64) (*FeeEntry, error) {
	args := q.Called(feeType, asset)
	return args.Get(0).(*FeeEntry), args.Error(1)
}

func (q *CoreQMock) AssetPairs() AssetPairsQ {
	args := q.Called()
	return args.Get(0).(AssetPairsQ)
}

func (q *CoreQMock) Accounts() AccountQI {
	args := q.Called()
	return args.Get(0).(AccountQI)
}

func (q *CoreQMock) Assets() AssetQI {
	args := q.Called()
	return args.Get(0).(AssetQI)
}

func (q *CoreQMock) GetRepo() *pgdb.DB {
	args := q.Called()
	return args.Get(0).(*pgdb.DB)
}

func (q *CoreQMock) AccountByAddress(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) ExchangeName(addy string) (*string, error) {
	args := q.Called(addy)
	return args.Get(0).(*string), args.Error(1)
}

func (q *CoreQMock) LedgerHeaderBySequence(dest interface{}, seq int32) error {
	args := q.Called(dest, seq)
	return args.Error(0)
}
func (q *CoreQMock) ElderLedger(dest *int32) error {
	args := q.Called(dest)
	return args.Error(0)
}
func (q *CoreQMock) LatestLedger(dest interface{}) error {
	args := q.Called(dest)
	return args.Error(0)
}

func (q *CoreQMock) SignersByAddress(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) PoliciesByExchangeID(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) TransactionByHash(dest interface{}, hash string) error {
	args := q.Called(dest, hash)
	return args.Error(0)
}
func (q *CoreQMock) TransactionsByLedger(dest interface{}, seq int32) error {
	args := q.Called(dest, seq)
	return args.Error(0)
}

func (q *CoreQMock) TransactionFeesByLedger(dest interface{}, seq int32) error {
	args := q.Called(dest, seq)
	return args.Error(0)
}

func (q *CoreQMock) FeeEntries() FeeEntryQI {
	args := q.Called()
	return args.Get(0).(FeeEntryQI)
}

func (q *CoreQMock) Query(query sq.Sqlizer) (*sqlx.Rows, error) {
	args := q.Called(query)
	return args.Get(0).(*sqlx.Rows), args.Error(1)
}

func (q *CoreQMock) FeeByOperationType(dest interface{}, operationType int) error {
	args := q.Called()
	return args.Error(0)
}
