package core

import (
	"database/sql"
	"encoding/hex"
	"fmt"

	"gitlab.com/tokend/go/hash"

	sq "github.com/Masterminds/squirrel"
)

var selectFees = sq.Select("f.fee_type", "f.asset", "f.subtype", "f.fixed", "f.percent", "f.lastmodified,"+
	"f.account_id, f.account_role, f.lower_bound, f.upper_bound, f.hash").
	From("fee_state f")

type FeeEntry struct {
	FeeType     int    `db:"fee_type"`
	Asset       string `db:"asset"`
	Fixed       int64  `db:"fixed"`
	Percent     int64  `db:"percent"`
	Subtype     int64  `db:"subtype"`
	AccountID   string `db:"account_id"`
	AccountRole int64  `db:"account_role"`
	LowerBound  int64  `db:"lower_bound"`
	UpperBound  int64  `db:"upper_bound"`
	Hash        string `db:"hash"`

	LastModified int32 `db:"lastmodified"`
}

// Fees loads all row from `fee_state`
func (q *Q) Fees(dest interface{}) error {
	sql := selectFees
	return q.Get(dest, sql)
}

func getFeesSelector(feeType int, asset string, subtype int64, account *Account) sq.SelectBuilder {
	query := selectFees
	baseString := fmt.Sprintf("type:%vasset:%vsubtype:%v", feeType, asset, subtype)
	filter := "hash IN ("
	orderFormat := "hash = '%s' DESC"
	if account != nil {
		hash1 := strHash(baseString + fmt.Sprintf("accountID:%v", account.AccountID))
		hash2 := strHash(baseString + fmt.Sprintf("accountType:%v", account.RoleID))
		filter += fmt.Sprintf("'%s', '%s', ", hash1, hash2)
		query = query.OrderBy(fmt.Sprintf(orderFormat, hash1), fmt.Sprintf(orderFormat, hash2))
	}

	hash3 := strHash(baseString)
	filter += fmt.Sprintf("'%s')", hash3)
	query = query.OrderBy(fmt.Sprintf(orderFormat, hash3))
	query = query.Where(filter)
	return query
}

func strHash(hashData string) string {
	rawHash := hash.Hash([]byte(hashData))
	return hex.EncodeToString(rawHash[:])
}

func (q *Q) FeeByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account, amount int64) (*FeeEntry, error) {
	var result FeeEntry

	query := getFeesSelector(feeType, asset, subtype, account).Where("lower_bound <= ? AND ? <= upper_bound", amount, amount).Limit(1)
	err := q.Get(&result, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (q *Q) FeesByTypeAssetAccount(feeType int, asset string, subtype int64, account *Account) ([]FeeEntry, error) {
	var result []FeeEntry
	query := getFeesSelector(feeType, asset, subtype, account).OrderBy("lower_bound DESC")
	err := q.Select(&result, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

// FeeEntryQ is a helper struct to aid in configuring queries that loads
// slices of FeeEntry structs.
type FeeEntryQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type FeeEntryQI interface {
	Select() ([]FeeEntry, error)
	ForAccountType(accountType *uint64) FeeEntryQI
	ForAccount(account string) FeeEntryQI
}

// FeeEntries provides a helper to filter the operations table with pre-defined
// filters.  See `FeeEntryQ` for the available filters.
func (q *Q) FeeEntries() FeeEntryQI {
	return &FeeEntryQ{
		parent: q,
		sql:    selectFees,
	}
}

func (q *FeeEntryQ) ForAccountType(accountType *uint64) FeeEntryQI {
	if q.Err != nil {
		return q
	}

	if accountType == nil {
		q.sql = q.sql.Where("account_role = ?", 0)
		return q
	}

	q.sql = q.sql.Where("account_role = ?", accountType)
	return q
}

func (q *FeeEntryQ) ForAccount(account string) FeeEntryQI {
	if q.Err != nil {
		return q
	}

	if account == "" {
		q.sql = q.sql.Where("(account_id = '' OR account_id IS NULL)")
		return q
	}

	q.sql = q.sql.Where("account_id = ?", account)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *FeeEntryQ) Select() (res []FeeEntry, err error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.Err = q.parent.Select(&res, q.sql)
	return res, q.Err
}
