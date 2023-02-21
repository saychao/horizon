package history

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2"
)

var selectAccount = sq.Select("ha.*").From("history_accounts ha")
var selectAccountsByOperation = sq.Select("ops.id as operation_id, acc.address").
	From("history_operations ops").
	LeftJoin("history_operation_participants opsparts ON opsparts.history_operation_id = ops.id").
	LeftJoin("history_accounts acc on acc.id = opsparts.history_account_id")

// AccountsQ is a helper struct to aid in configuring queries that loads
// slices of account structs.
type AccountsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type AccountsQI interface {
	AccountsByOperationIDs(operationIDs []int64) (map[int64][]Account, error)
	Page(page db2.PageQuery) AccountsQI
	Select(dest interface{}) error
}

// Accounts provides a helper to filter rows from the `history_accounts` table
// with pre-defined filters.  See `AccountsQ` methods for the available filters.
func (q *Q) Accounts() AccountsQI {
	return &AccountsQ{
		parent: q,
		sql:    selectAccount,
	}
}

// AccountByAddress loads a row from `history_accounts`, by address
func (q *Q) AccountByAddress(dest interface{}, addy string) error {
	sql := selectAccount.Limit(1).Where("ha.address = ?", addy)
	return q.Get(dest, sql)
}

// AccountByID loads a row from `history_accounts`, by id
func (q *Q) AccountByID(dest interface{}, id int64) error {
	sql := selectAccount.Limit(1).Where("ha.id = ?", id)
	return q.Get(dest, sql)
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *AccountsQ) Page(page db2.PageQuery) AccountsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ha.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *AccountsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

type OperationAccount struct {
	OperationID int64 `db:"operation_id"`
	Account
}

func (q *AccountsQ) AccountsByOperationIDs(operationIDs []int64) (map[int64][]Account, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	if len(operationIDs) == 0 {
		return map[int64][]Account{}, nil
	}

	q.sql = selectAccountsByOperation.Where(sq.Eq{"ops.id": operationIDs})
	var operationUsers []OperationAccount
	err := q.parent.Select(&operationUsers, q.sql)
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]Account)
	for i := range operationUsers {
		operationAccount := operationUsers[i]
		if result[operationAccount.OperationID] == nil {
			result[operationAccount.OperationID] = make([]Account, 0, 1)
		}

		result[operationAccount.OperationID] = append(result[operationAccount.OperationID], operationAccount.Account)
	}

	return result, nil
}
