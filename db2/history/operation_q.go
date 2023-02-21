package history

import (
	"github.com/saychao/horizon/db2"
	"time"

	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2/sqx"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

var selectOperation = sq.Select("distinct on (ho.id) ho.*").
	From("history_operations ho").
	LeftJoin("history_operation_participants hop on hop.history_operation_id = ho.id")

// OperationsQ is a helper struct to aid in configuring queries that loads
// slices of Operation structs.
type OperationsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type OperationsQI interface {
	ForAccount(aid string) OperationsQI
	ForAccountType(accountType int32) OperationsQI
	ForBalance(bid string) OperationsQI
	ForAsset(asset string) OperationsQI
	ForTypes(opTypes []xdr.OperationType) OperationsQI
	ForTx(txhash string) OperationsQI
	ForReference(reference string) OperationsQI
	Since(ts time.Time) OperationsQI
	To(ts time.Time) OperationsQI
	CompletedOnly() OperationsQI
	PendingOnly() OperationsQI
	// required prior to exchange and asset filtering
	JoinOnBalance() OperationsQI
	// JoinOnAccount required to filter on account ID and account type
	JoinOnAccount() OperationsQI
	Page(page db2.PageQuery) OperationsQI
	Select(dest interface{}) error

	// Manage Offer
	// WithoutCancelingManagerOffer - don't load manage offer operations which cancel offer
	WithoutCancelingManagerOffer() OperationsQI
	// WithoutCanceled - filters canceled operations
	WithoutCanceled() OperationsQI

	// WithoutExternallyFullyMatched = don't load manage offer operations with ExternallyFullyMatched state
	WithoutExternallyFullyMatched() OperationsQI

	// OnlyReviewOfAtomicSwapRequest - load review request ops only with atomic swap request type
	OnlyReviewOfAtomicSwapRequest() OperationsQI

	Participants(dest map[int64]*OperationParticipants) error
}

// Operations provides a helper to filter the operations table with pre-defined
// filters.  See `OperationsQ` for the available filters.
func (q *Q) Operations() OperationsQI {
	return &OperationsQ{
		parent: q,
		sql:    selectOperation,
	}
}

// OperationByID loads a single operation with `id` into `dest`
func (q *Q) OperationByID(dest interface{}, id int64) error {
	sql := selectOperation.
		Limit(1).
		Where("ho.id = ?", id)

	return q.Get(dest, sql)
}

func (q *OperationsQ) Participants(dest map[int64]*OperationParticipants) error {
	ids := []int64{}
	for opid, _ := range dest {
		ids = append(ids, opid)
	}

	if len(ids) == 0 {
		return nil
	}

	stmt := sq.Select("hop.history_operation_id, ha.address as account_id, hop.balance_id, hop.effects").
		From("history_operation_participants hop").
		LeftJoin("history_accounts ha on ha.id = hop.history_account_id"). // join to get account addresses
		Where(sq.Eq{"hop.history_operation_id": ids})

	var participants []*Participant
	err := q.parent.Select(&participants, stmt)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		opid := participant.OperationID
		dest[opid].Participants = append(dest[opid].Participants, participant)
	}

	return nil
}

func (q *OperationsQ) JoinOnAccount() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Join("history_accounts ha on ha.id = hop.history_account_id")

	return q
}

// ForAccount filters the operations collection to a specific account
func (q *OperationsQ) ForAccount(address string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ha.address = ?", address)

	return q
}

func (q *OperationsQ) ForAccountType(accountType int32) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ha.account_type = ?", accountType)

	return q
}

func (q *OperationsQ) JoinOnBalance() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.LeftJoin("history_balances hb on hb.balance_id=hop.balance_id")

	return q
}

func (q *OperationsQ) ForAsset(asset string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hb.asset = ?", asset)

	return q
}

func (q *OperationsQ) ForBalance(balanceID string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hop.balance_id = ?", balanceID)

	return q
}

func (q *OperationsQ) ForTypes(opTypes []xdr.OperationType) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"ho.type": opTypes})

	return q
}

func (q *OperationsQ) ForTx(txhash string) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Join("history_transactions ht on ht.id=ho.transaction_id").
		Where("ht.transaction_hash = ?", txhash)

	return q
}

func (q *OperationsQ) ForReference(reference string) OperationsQI {
	if q.Err != nil {
		return q
	}

	// FIXME might(will) not work for all operation types, works at least for payments and issuances
	q.sql = q.sql.Where("ho.details->>'reference' ilike ?", fmt.Sprintf("%%%s%%", reference))

	return q
}

func (q *OperationsQ) Since(ts time.Time) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.ledger_close_time >= ?", ts)

	return q
}

func (q *OperationsQ) To(ts time.Time) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.ledger_close_time <= ?", ts)

	return q
}

func (q *OperationsQ) CompletedOnly() OperationsQI {
	if q.Err != nil {
		return q
	}

	query, values := sqx.In("state",
		OperationStateSuccess,
		OperationStateRejected,
		OperationStateCanceled,
		OperationStateFailed,
		OperationStateFullyMatched)

	q.sql = q.sql.Where(query, values...)
	return q
}

func (q *OperationsQ) PendingOnly() OperationsQI {
	if q.Err != nil {
		return q
	}

	query, values := sqx.In("state",
		OperationStatePending,
		OperationStatePartiallyMatched)

	q.sql = q.sql.Where(query, values...)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *OperationsQ) Page(page db2.PageQuery) OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ho.id")
	return q
}

func (q *OperationsQ) WithoutCancelingManagerOffer() OperationsQI {
	if q.Err != nil {
		return q
	}

	// 'amount' field in 'details' jsonb has type string, thus required to pass amount.String to query
	q.sql = q.sql.Where("(ho.type <> ? OR (ho.details->>'amount') != ?)", xdr.OperationTypeManageOffer,
		amount.String(0))
	return q
}

func (q *OperationsQ) WithoutExternallyFullyMatched() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.state <> ?", OperationStateExternallyFullyMatched)
	return q
}

func (q *OperationsQ) OnlyReviewOfAtomicSwapRequest() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.request_type IS NULL OR ho.request_type == ?",
		xdr.ReviewableRequestTypeCreateAtomicSwapAsk.ShortString())
	return q
}

// WithoutCanceled - filters canceled operations
func (q *OperationsQ) WithoutCanceled() OperationsQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ho.state <> ?", OperationStateCanceled)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OperationsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
