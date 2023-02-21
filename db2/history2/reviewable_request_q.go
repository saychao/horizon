package history2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ReviewableRequestsQ - helper struct to get reviewable requests from db
type ReviewableRequestsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewReviewableRequestsQ - creates new instance of ReviewableRequestsQ
func NewReviewableRequestsQ(repo *pgdb.DB) ReviewableRequestsQ {
	return ReviewableRequestsQ{
		repo: repo,
		selector: sq.Select(
			"reviewable_requests.id",
			"reviewable_requests.requestor",
			"reviewable_requests.reviewer",
			"reviewable_requests.reference",
			"reviewable_requests.reject_reason",
			"reviewable_requests.request_type",
			"reviewable_requests.request_state",
			"reviewable_requests.hash",
			"reviewable_requests.details",
			"reviewable_requests.created_at",
			"reviewable_requests.updated_at",
			"reviewable_requests.all_tasks",
			"reviewable_requests.pending_tasks",
			"reviewable_requests.external_details",
		).From("reviewable_requests"),
	}
}

func (q ReviewableRequestsQ) FilterByCreatedAtAfter(tsUTC int64) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.created_at at time zone 'UTC' >= to_timestamp(?::numeric)", tsUTC)
	return q
}

func (q ReviewableRequestsQ) FilterByCreatedAtBefore(tsUTC int64) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.created_at at time zone 'UTC' < to_timestamp(?::numeric)", tsUTC)
	return q
}

// FilterByState - returns q with filter by request state
func (q ReviewableRequestsQ) FilterByState(state uint64) ReviewableRequestsQ { // TODO: move state to regources?
	q.selector = q.selector.Where("reviewable_requests.request_state = ?", state)
	return q
}

// FilterByRequestorAddress - returns q with filter by requestor address
func (q ReviewableRequestsQ) FilterByRequestorAddress(address string) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.requestor = ?", address)
	return q
}

// FilterByReviewerAddress - returns q with filter by reviewer
func (q ReviewableRequestsQ) FilterByReviewerAddress(address string) ReviewableRequestsQ {
	q.selector = q.selector.Where(sq.Eq{"reviewable_requests.reviewer": address})
	return q
}

// FilterByPendingTasks - returns q with filter by pending tasks
func (q ReviewableRequestsQ) FilterByPendingTasks(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("pending_tasks & ? = ?", mask, mask)
	return q
}

// FilterByPendingTasksAnyOf - returns q with filter by any of passed pending tasks
func (q ReviewableRequestsQ) FilterByPendingTasksAnyOf(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("pending_tasks & ? <> 0", mask)
	return q
}

// FilterPendingTasksNotSet - returns q with filter by pending tasks that aren't set
func (q ReviewableRequestsQ) FilterPendingTasksNotSet(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("~pending_tasks & ? = ?", mask, mask)
	return q
}

// FilterByAllTasks - returns q with filter by all tasks
func (q ReviewableRequestsQ) FilterByAllTasks(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("all_tasks & ? = ?", mask, mask)
	return q
}

// FilterByAllTasksAnyOf - returns q with filter by any of passed all tasks
func (q ReviewableRequestsQ) FilterByAllTasksAnyOf(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("all_tasks & ? <> 0", mask)
	return q
}

// FilterByAllTasksNotSet - returns q with filter by all tasks that aren't set
func (q ReviewableRequestsQ) FilterByAllTasksNotSet(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("~all_tasks & ? = ?", mask, mask)
	return q
}

// FilterByRequestType - returns q with filter by request type
func (q ReviewableRequestsQ) FilterByRequestType(requestType uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("request_type = ?", requestType)
	return q
}

// FilterByID - returns q with filter by ID
func (q ReviewableRequestsQ) FilterByID(id uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.id = ?", id)
	return q
}

func (q ReviewableRequestsQ) FilterByAssetCreateAsset(asset string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_asset,asset}' = ?", asset)
	return q
}

func (q ReviewableRequestsQ) FilterByAssetUpdateAsset(asset string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{update_asset,asset}' = ?", asset)
	return q
}

func (q ReviewableRequestsQ) FilterByCreatePreIssuanceAsset(asset string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_pre_issuance,asset}' = ?", asset)
	return q
}

func (q ReviewableRequestsQ) FilterByCreateIssuanceAsset(asset ...string) ReviewableRequestsQ {
	q.selector = q.selector.Where(sq.Eq{"details#>>'{create_issuance,asset}'": asset})
	return q
}

func (q ReviewableRequestsQ) FilterByCreateIssuanceReceiver(receiver string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_issuance,receiver}' = ?", receiver)
	return q
}

func (q ReviewableRequestsQ) FilterByWithdrawBalance(balance string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_withdraw,balance_id}' = ?", balance)
	return q
}

func (q ReviewableRequestsQ) FilterByWithdrawAssets(assets ...string) ReviewableRequestsQ {
	q.selector = q.selector.Where(sq.Eq{"details#>>'{create_withdraw,asset}'": assets})
	return q
}

func (q ReviewableRequestsQ) FilterByAmlAlertBalance(balance string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_aml_alert,balance_id}' = ?", balance)
	return q
}

func (q ReviewableRequestsQ) FilterByChangeRoleAccount(account string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{change_role,destination_account}' = ?", account)
	return q
}

func (q ReviewableRequestsQ) FilterByChangeRoleToSet(accountRole int32) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{change_role,account_role_to_set}' = ?", accountRole)
	return q
}

func (q ReviewableRequestsQ) FilterByCreateAtomicSwapAskBalance(balance string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_atomic_swap_ask,base_balance}' = ?", balance)
	return q
}

func (q ReviewableRequestsQ) FilterByAtomicSwapQuoteAsset(code string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_atomic_swap_bid,quote_asset}' = ?", code)
	return q
}

func (q ReviewableRequestsQ) FilterByAtomicSwapAskID(id uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_atomic_swap_bid,bid_id}' = ?", id) //bid_id because there is mistake in reviewable request details
	return q
}

func (q ReviewableRequestsQ) FilterByAtomicSwapAskIDs(ids []uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where(sq.Eq{"details#>>'{create_atomic_swap_bid,bid_id}'": ids}) //bid_id because there is mistake in reviewable request details
	return q
}

func (q ReviewableRequestsQ) FilterBySaleBaseAsset(asset string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_sale,base_asset}' = ?", asset)
	return q
}

func (q ReviewableRequestsQ) FilterBySaleQuoteAsset(asset string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{sale,quote_asset}' = ?", asset)
	return q
}

func (q ReviewableRequestsQ) FilterByCreatePollPermissionType(permissionType uint32) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_poll,permission_type}' = ?", permissionType)
	return q
}

func (q ReviewableRequestsQ) FilterByDataCreationSecurityType(securityType uint32) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{data_creation,security_type}' = ?", securityType)
	return q
}

func (q ReviewableRequestsQ) FilterByCreateDeferredPaymentDestination(destination string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_deferred_payment,destination_account}' = ?", destination)
	return q
}

func (q ReviewableRequestsQ) FilterByCreatePollVoteConfirmationRequired(voteConfirmationRequired bool) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_poll,vote_confirmation_required}' = ?", voteConfirmationRequired)
	return q
}

func (q ReviewableRequestsQ) FilterByCreatePollResultProvider(resultProviderID string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_poll,result_provider_id}' = ?", resultProviderID)
	return q
}

func (q ReviewableRequestsQ) FilterByKYCRecoveryTargetAccount(address string) ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{kyc_recovery,account}' = ?", address)
	return q
}

func (q ReviewableRequestsQ) FilterByRegistrationType(adminRegistered bool) ReviewableRequestsQ {
	if adminRegistered {
		q.selector = q.selector.Where("details#>>'{change_role,creator_details,admin_registered}' = ?", true)
		return q
	}
	q.selector = q.selector.Where("(details#>>'{change_role,creator_details,admin_registered}')::boolean IS NOT TRUE")
	return q
}

func (q ReviewableRequestsQ) FilterByTimestampGT(timestamp int64) ReviewableRequestsQ {
	q.selector = q.selector.Where(sq.GtOrEq{"extract(epoch from updated_at)": timestamp})
	return q
}

func (q ReviewableRequestsQ) FilterWhoNotNull() ReviewableRequestsQ {
	q.selector = q.selector.Where("details#>>'{create_issuance,external_details,who}' is not null")
	return q
}

// GetByID loads a row from `reviewable_requests`, by ID
// returns nil, nil - if request does not exists
func (q ReviewableRequestsQ) GetByID(id uint64) (*ReviewableRequest, error) {
	return q.FilterByID(id).Get()
}

// Page - apply paging params to the query
func (q ReviewableRequestsQ) Page(pageParams pgdb.CursorPageParams) ReviewableRequestsQ {
	q.selector = pageParams.ApplyTo(q.selector, "reviewable_requests.id")
	return q
}

// Get - loads a row from `reviewable_requests`
// returns nil, nil - if request does not exists
// returns error if more than one ReviewableRequest found
func (q ReviewableRequestsQ) Get() (*ReviewableRequest, error) {
	var result ReviewableRequest
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}

// Select - selects ReviewableRequest from db using specified filters. Returns nil, nil - if one does not exists
func (q ReviewableRequestsQ) Select() ([]ReviewableRequest, error) {
	var result []ReviewableRequest
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select reviewable requests")
	}

	return result, nil
}
