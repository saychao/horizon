package history

import (
	"database/sql"
	"github.com/SafeRE-IT/horizon/db2"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2/sqx"
	"gitlab.com/tokend/go/xdr"
)

// ReviewableRequestQI - provides methods to operate reviewable request
type ReviewableRequestQI interface {
	// Insert - inserts new request
	Insert(request ReviewableRequest) error
	// Update - update request using it's ID
	Update(request ReviewableRequest) error
	// UpdateStates - update state of requests
	UpdateStates(requestIDs []int64, state ReviewableRequestState) error
	// ByID - selects request by id. Returns nil, nil if not found
	ByID(requestID uint64) (*ReviewableRequest, error)
	// ByID - selects request by id. Returns nil, nil if not found
	ByIDs(requestIDs []int64) ReviewableRequestQI
	// Cancel - sets request state to `ReviewableRequestStateCanceled`
	Cancel(requestID uint64) error
	// Approve - sets request state to ReviewableRequestStateApproved and cleans reject reason
	Approve(requestID uint64) error
	// PermanentReject - sets request state to ReviewableRequestStatePermanentlyRejected and sets reject reason
	PermanentReject(requestID uint64, rejectReason string) error
	// ForRequestor - filters requests by requestor
	ForRequestor(requestor string) ReviewableRequestQI
	// ForReviewer - filters requests by reviewer
	ForReviewer(reviewer string) ReviewableRequestQI
	// ForCounterparty - filters requests by reviewer or requestor
	ForCounterparty(counterparty string) ReviewableRequestQI
	// ForState - filters requests by state
	ForState(state int64) ReviewableRequestQI
	// ForType - filters requests by type
	ForType(requestType int64) ReviewableRequestQI
	// ForRoles - filters requests by request type
	ForTypes(requestTypes []xdr.ReviewableRequestType) ReviewableRequestQI
	// Page specifies the paging constraints for the query being built by `q`.
	Page(page db2.PageQuery) ReviewableRequestQI
	// UpdatedAfter - selects requests updated after given timestamp
	UpdatedAfter(timestamp int64) ReviewableRequestQI
	// Select loads the results of the query specified by `q`
	Select() ([]ReviewableRequest, error)
	// Count loads count of the results of the query specified by `q`
	Count() (int64, error)
	CountQuery() ReviewableRequestQI

	// Request Type specific filters. Filter for request type must be applied separately

	// Asset Management
	// AssetManagementByAsset - filters asset management requests by asset
	AssetManagementByAsset(assetCode string) ReviewableRequestQI

	// PreIssuance
	// PreIssuanceByAsset - filters pre issuance requests by asset
	PreIssuanceByAsset(assetCode string) ReviewableRequestQI

	// Issuance
	// IssuanceByAsset - filters issuance requests by asset
	IssuanceByAsset(assetCode string) ReviewableRequestQI

	// Withdraw
	// WithdrawalByDestAsset - filters withdrawal requests by dest asset
	WithdrawalByDestAsset(assetCode string) ReviewableRequestQI

	// Sales
	// SalesByBaseAsset - filters sale requests by base asset
	SalesByBaseAsset(assetCode string) ReviewableRequestQI

	// Limits
	// LimitsByDocHash - filters limits request by document hash
	LimitsByDocHash(hash string) ReviewableRequestQI

	// Contracts
	// ContractsByContractNumber - filters contract requests by contract number
	ContractsByContractNumber(contractNumber string) ReviewableRequestQI
	// ContractsByStartTime - filters contract requests by contract start time
	ContractsByStartTime(seconds int64) ReviewableRequestQI
	// ContractsByEndTime - filters contract requests by contract start time
	ContractsByEndTime(seconds int64) ReviewableRequestQI

	// Invoices
	// InvoicesByContract - filters invoice requests by contract id
	InvoicesByContract(contractID int64) ReviewableRequestQI
	// UpdateInvoicesStates - update state of invoice requests by contract id
	UpdateInvoicesStates(state ReviewableRequestState, oldStates []ReviewableRequestState, contractID int64) error

	// KYC
	// KYCByAccountToUpdateKYC - filters update KYC requests by accountID of the owner of KYC
	KYCByAccountToUpdateKYC(accountID string) ReviewableRequestQI
	// KYCByMaskSet - filters update KYC requests by mask which must be set. If maskSetPartialEq is false, request will be returned
	// even if only part of the mask is set
	ByMaskSet(mask int64, maskSetPartialEq bool) ReviewableRequestQI
	// KYCByMaskNotSet - filters update KYC requests by mask which must not be set
	ByMaskNotSet(mask int64) ReviewableRequestQI
	// KYCByAccountTypeToSet - filters update KYC requests by account type which must be set.
	KYCByAccountTypeToSet(accountTypeToSet xdr.Uint64) ReviewableRequestQI

	// Atomic swap
	// AtomicSwapByBidID - filters atomic swap requests by bid id
	AtomicSwapByBidID(bidID int64) ReviewableRequestQI
}

type ReviewableRequestQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// Insert - inserts new request
func (q *ReviewableRequestQ) Insert(request ReviewableRequest) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Insert("reviewable_request").SetMap(map[string]interface{}{
		"id":               request.ID,
		"requestor":        request.Requestor,
		"reviewer":         request.Reviewer,
		"reference":        request.Reference,
		"reject_reason":    request.RejectReason,
		"request_type":     request.RequestType,
		"request_state":    request.RequestState,
		"hash":             request.Hash,
		"details":          request.Details,
		"created_at":       request.CreatedAt,
		"updated_at":       request.UpdatedAt,
		"all_tasks":        request.AllTasks,
		"pending_tasks":    request.PendingTasks,
		"external_details": request.ExternalDetails,
	})

	err := q.parent.Exec(query)
	return err
}

// Update - update request using it's ID
func (q *ReviewableRequestQ) Update(request ReviewableRequest) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").SetMap(map[string]interface{}{
		"requestor":        request.Requestor,
		"reviewer":         request.Reviewer,
		"reject_reason":    request.RejectReason,
		"request_type":     request.RequestType,
		"request_state":    request.RequestState,
		"hash":             request.Hash,
		"details":          request.Details,
		"updated_at":       request.UpdatedAt,
		"all_tasks":        request.AllTasks,
		"pending_tasks":    request.PendingTasks,
		"external_details": request.ExternalDetails,
	}).Where("id = ?", request.ID)

	err := q.parent.Exec(query)
	return err
}

func (q *ReviewableRequestQ) UpdateStates(requestIDs []int64, state ReviewableRequestState) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").
		Set("request_state", state).
		Where(sq.Eq{"id": requestIDs})

	err := q.parent.Exec(query)
	return err
}

// ByID - selects request by id. Returns nil, nil if not found
func (q *ReviewableRequestQ) ByID(requestID uint64) (*ReviewableRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := q.sql.Where("rr.id = ?", requestID)

	var result ReviewableRequest
	err := q.parent.Get(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (q *ReviewableRequestQ) ByIDs(requestIDs []int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"rr.id": requestIDs})
	return q
}

// Cancel - sets request state to `ReviewableRequestStateCanceled`
func (q *ReviewableRequestQ) Cancel(requestID uint64) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").
		Set("request_state", ReviewableRequestStateCanceled).Where("id = ?", requestID)

	err := q.parent.Exec(query)
	return err
}

// Approve - sets request state to ReviewableRequestStateApproved and cleans reject reason
func (q *ReviewableRequestQ) Approve(requestID uint64) error {
	return q.setStateRejectReason(requestID, ReviewableRequestStateApproved, "")
}

// PermanentReject - sets request state to ReviewableRequestStatePermanentlyRejected and sets reject reason
func (q *ReviewableRequestQ) PermanentReject(requestID uint64, rejectReason string) error {
	return q.setStateRejectReason(requestID, ReviewableRequestStatePermanentlyRejected, rejectReason)
}

func (q *ReviewableRequestQ) setStateRejectReason(requestID uint64, requestState ReviewableRequestState, rejectReason string) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").
		Set("request_state", requestState).
		Set("reject_reason", rejectReason).Where("id = ?", requestID)

	err := q.parent.Exec(query)
	return err
}

// ForRequestor - filters requests by requestor
func (q *ReviewableRequestQ) ForRequestor(requestor string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.requestor = ?", requestor)
	return q
}

// ForReviewer - filters requests by reviewer
func (q *ReviewableRequestQ) ForReviewer(reviewer string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.reviewer = ?", reviewer)
	return q
}

func (q *ReviewableRequestQ) ForCounterparty(counterparty string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("((rr.reviewer = ?) or (rr.requestor = ?))", counterparty, counterparty)
	return q
}

// ForState - filters requests by state
func (q *ReviewableRequestQ) ForState(state int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.request_state = ?", state)
	return q
}

// ForType - filters requests by type
func (q *ReviewableRequestQ) ForType(requestType int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.request_type = ?", requestType)
	return q
}

// ForRoles - filters requests by request type
func (q *ReviewableRequestQ) ForTypes(requestTypes []xdr.ReviewableRequestType) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	if len(requestTypes) == 0 {
		return q
	}

	query, values := sqx.InForReviewableRequestTypes("rr.request_type", requestTypes...)

	q.sql = q.sql.Where(query, values...)
	return q
}

// UpdatedAfter - selects requests updated after given timestamp
func (q *ReviewableRequestQ) UpdatedAfter(timestamp int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	tm := time.Unix(timestamp, 0)
	tmf := tm.Format(time.RFC3339)

	q.sql = q.sql.Where("rr.updated_at >= ?::timestamp", tmf)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *ReviewableRequestQ) Page(page db2.PageQuery) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "rr.id")
	return q
}

// Select loads the results of the query specified by `q`
func (q *ReviewableRequestQ) Select() ([]ReviewableRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []ReviewableRequest
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

func (q *ReviewableRequestQ) CountQuery() ReviewableRequestQI {
	if q.Err != nil {
		return q
	}
	q.sql = sq.Select("COUNT(*)").From("reviewable_request rr")
	return q
}

func (q *ReviewableRequestQ) Count() (int64, error) {
	if q.Err != nil {
		return 0, q.Err
	}

	var result int64
	q.Err = q.parent.Get(&result, q.sql)
	return result, q.Err
}

func (q *ReviewableRequestQ) AssetManagementByAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(rr.details->'create_asset'->>'asset' = ? OR rr.details->'update_asset'->>'asset' = ?)", assetCode, assetCode)
	return q
}

// PreIssuance
// PreIssuanceByAsset - filters pre issuance requests by asset
func (q *ReviewableRequestQ) PreIssuanceByAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'create_pre_issuance'->>'asset' = ?", assetCode)
	return q
}

// Issuance
// IssuanceByAsset - filters issuance requests by asset
func (q *ReviewableRequestQ) IssuanceByAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'create_issuance'->>'asset' = ?", assetCode)
	return q
}

// Withdraw
// WithdrawalByDestAsset - filters withdrawal requests by dest asset
func (q *ReviewableRequestQ) WithdrawalByDestAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Join("history_balances hb on details #>> '{create_withdraw, balance_id}' = (hb.balance_id) OR details #>> '{two_step_withdrawal, balance_id}' = (hb.balance_id)").Where("hb.asset = ?", assetCode)
	return q
}

// Sales
// SalesByBaseAsset - filters sale requests by base asset
func (q *ReviewableRequestQ) SalesByBaseAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'create_sale'->>'base_asset' = ?", assetCode)
	return q
}

// Limits
// LimitsByDocHash - filters limits request by document hash
func (q *ReviewableRequestQ) LimitsByDocHash(hash string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'update_limits'->>'document_hash' = ?", hash)
	return q
}

func (q *ReviewableRequestQ) ContractsByContractNumber(contractNumber string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'contract'->'details'->>'contract_number' = ?", contractNumber)
	return q
}

func (q *ReviewableRequestQ) ContractsByStartTime(seconds int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'contract'->>'start_time' >= ?", time.Unix(seconds, 0).UTC())
	return q
}

func (q *ReviewableRequestQ) ContractsByEndTime(seconds int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'contract'->>'end_time' >= ?", time.Unix(seconds, 0).UTC())
	return q
}

func (q *ReviewableRequestQ) InvoicesByContract(contractID int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'invoice'->>'contract_id' = ?", contractID)
	return q
}

func (q *ReviewableRequestQ) UpdateInvoicesStates(state ReviewableRequestState,
	oldStates []ReviewableRequestState,
	contractID int64,
) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").
		Set("request_state", state).
		Where("details->'invoice'->>'contract_id' = ?", contractID).
		Where(sq.Eq{"request_state": oldStates})

	err := q.parent.Exec(query)
	return err
}

func (q *ReviewableRequestQ) AtomicSwapByBidID(bidID int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("rr.details->'create_atomic_swap'->>'bid_id' = ?", bidID)
	return q
}

// KYC
// KYCByAccountToUpdateKYC - filters update KYC requests by accountID of the owner of KYC
func (q *ReviewableRequestQ) KYCByAccountToUpdateKYC(accountID string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'change_role'->>'destination_account' = ?", accountID)
	return q
}

// ByMaskSet - filters update KYC requests by mask which must be set. If mustBeEq is false, request will be returned
// even if only part of the mask is set
func (q *ReviewableRequestQ) ByMaskSet(mask int64, maskSetPartialEq bool) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	if maskSetPartialEq {
		q.sql = q.sql.Where("rr.pending_tasks & ? <> 0", mask)
	} else {
		q.sql = q.sql.Where("rr.pending_tasks & ? = ?", mask, mask)
	}
	return q
}

// ByMaskNotSet - filters update KYC requests by mask which must not be set
func (q *ReviewableRequestQ) ByMaskNotSet(mask int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("~rr.pending_tasks & ? = ?", mask, mask)
	return q
}

// KYCByAccountTypeToSet - filters update KYC requests by account type which must be set.
func (q *ReviewableRequestQ) KYCByAccountTypeToSet(accountTypeToSet xdr.Uint64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(rr.details->'change_role'->'account_role_to_set'->>'int')::integer = ?", int32(accountTypeToSet))
	return q
}

var selectReviewableRequest = sq.Select("rr.id", "rr.requestor", "rr.reviewer", "rr.reference", "rr.reject_reason",
	"rr.request_type", "rr.request_state", "rr.hash", "rr.details", "rr.created_at", "rr.updated_at", "rr.all_tasks", "rr.pending_tasks",
	"rr.external_details").From("reviewable_request rr")
