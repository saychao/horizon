package history

import (
	"github.com/saychao/horizon/db2"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type QMock struct {
	mock.Mock
}

type AccountsQMock struct {
	mock.Mock
}

type LedgersQMock struct {
	mock.Mock
}

type OperationsQMock struct {
	mock.Mock
}

type TransactionsQMock struct {
	mock.Mock
}

func (q *QMock) GetRepo() *pgdb.DB {
	return nil
}

func (q *QMock) ElderLedger(dest interface{}) error {
	return q.Called(dest).Error(0)
}
func (q *QMock) LatestLedger(dest interface{}) error {
	return q.Called(dest).Error(0)
}
func (q *QMock) OldestOutdatedLedgers(dest interface{}, currentVersion int) error {
	return q.Called(dest, currentVersion).Error(0)
}

// Accounts
func (q *QMock) Accounts() AccountsQI {
	return q.Called().Get(0).(AccountsQI)
}
func (q *QMock) AccountByAddress(dest interface{}, addy string) error {
	return q.Called(dest, addy).Error(0)
}
func (q *QMock) AccountByID(dest interface{}, id int64) error {
	return q.Called(dest, id).Error(0)
}
func (q *AccountsQMock) Page(page db2.PageQuery) AccountsQI {
	return q.Called(page).Get(0).(AccountsQI)
}
func (q *AccountsQMock) Select(dest interface{}) error {
	return q.Called(dest).Error(0)
}

// Ledgers
func (q *QMock) Ledgers() LedgersQI {
	return q.Called().Get(0).(LedgersQI)
}
func (q *QMock) LedgerBySequence(dest interface{}, seq int32) error {
	return q.Called(dest, seq).Error(0)
}
func (q *LedgersQMock) Page(page db2.PageQuery) LedgersQI {
	return q.Called(page).Get(0).(LedgersQI)
}
func (q *LedgersQMock) Select(dest interface{}) error {
	return q.Called(dest).Error(0)
}

// Operations
func (q *QMock) Operations() OperationsQI {
	return q.Called().Get(0).(OperationsQI)
}
func (q *QMock) OperationByID(dest interface{}, id int64) error {
	return q.Called(dest).Error(0)
}

func (q *OperationsQMock) ForAccount(aid string) OperationsQI {
	return q.Called(aid).Get(0).(OperationsQI)
}
func (q *OperationsQMock) ForLedger(seq int32) OperationsQI {
	return q.Called(seq).Get(0).(OperationsQI)
}
func (q *OperationsQMock) ForTransaction(hash string) OperationsQI {
	return q.Called(hash).Get(0).(OperationsQI)
}
func (q *OperationsQMock) OnlyPayments() OperationsQI {
	return q.Called().Get(0).(OperationsQI)
}
func (q *OperationsQMock) Page(page db2.PageQuery) OperationsQI {
	return q.Called(page).Get(0).(OperationsQI)
}
func (q *OperationsQMock) Select(dest interface{}) error {
	return q.Called(dest).Error(0)
}

// Transactions
func (q *QMock) Transactions() TransactionsQI {
	return q.Called().Get(0).(TransactionsQI)
}
func (q *QMock) TransactionByHash(dest interface{}, hash string) error {
	return q.Called(dest).Error(0)
}

func (q *TransactionsQMock) ForAccount(aid string) TransactionsQI {
	return q.Called(aid).Get(0).(TransactionsQI)
}
func (q *TransactionsQMock) ForLedger(seq int32) TransactionsQI {
	return q.Called(seq).Get(0).(TransactionsQI)
}
func (q *TransactionsQMock) Page(page db2.PageQuery) TransactionsQI {
	return q.Called(page).Get(0).(TransactionsQI)
}
func (q *TransactionsQMock) Select(dest interface{}) error {
	return q.Called(dest).Error(0)
}
