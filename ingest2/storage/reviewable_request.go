package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ReviewableRequest is helper struct to operate with `reviewable requests`
type ReviewableRequest struct {
	repo *pgdb.DB
}

// NewReviewableRequest - creates new instance
func NewReviewableRequest(repo *pgdb.DB) *ReviewableRequest {
	return &ReviewableRequest{
		repo: repo,
	}
}

// Insert - inserts new request
func (q *ReviewableRequest) Insert(request history2.ReviewableRequest) error {
	query := sq.Insert("reviewable_requests").SetMap(map[string]interface{}{
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

	err := q.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert reviewable_request")
	}

	return nil
}

// Update - update request using it's ID
func (q *ReviewableRequest) Update(request history2.ReviewableRequest) error {
	query := sq.Update("reviewable_requests").SetMap(map[string]interface{}{
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

	err := q.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to do full update of reviewable_request",
			logan.F{"request_id": request.ID})
	}

	return nil
}

// Cancel - sets request state to `ReviewableRequestStateCanceled`
func (q *ReviewableRequest) Cancel(requestID uint64) error {
	return q.setStateRejectReason(requestID, history2.ReviewableRequestStateCanceled, nil)
}

// Approve - sets request state to ReviewableRequestStateApproved and cleans reject reason
func (q *ReviewableRequest) Approve(requestID uint64) error {
	return q.setStateRejectReason(requestID, history2.ReviewableRequestStateApproved, nil)
}

// PermanentReject - sets request state to ReviewableRequestStatePermanentlyRejected and sets reject reason
func (q *ReviewableRequest) PermanentReject(requestID uint64, rejectReason string) error {
	return q.setStateRejectReason(requestID, history2.ReviewableRequestStatePermanentlyRejected, &rejectReason)
}

func (q *ReviewableRequest) setStateRejectReason(requestID uint64, requestState history2.ReviewableRequestState,
	rejectReason *string) error {

	query := sq.Update("reviewable_requests").
		Set("request_state", requestState).Where("id = ?", requestID)
	if rejectReason != nil {
		query = query.Set("reject_reason", rejectReason)
	}

	err := q.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to set reject reason", logan.F{
			"request_id":    requestID,
			"request_state": requestState,
		})
	}

	return nil
}
