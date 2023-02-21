package reviewablerequest

import (
	"encoding/json"
	"github.com/saychao/horizon/db2/history"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
)

func PopulateReviewableRequest(request *history.ReviewableRequest) (r *regources.ReviewableRequest, err error) {
	r = &regources.ReviewableRequest{}

	r.ID = uint64(request.ID)
	r.PT = request.PagingToken()
	r.Requestor = request.Requestor
	r.Reviewer = request.Reviewer
	r.Reference = request.Reference
	r.RejectReason = request.RejectReason
	r.StateName = request.RequestState.String()
	r.State = int32(request.RequestState)
	r.Hash = request.Hash
	r.AllTasks = request.AllTasks
	r.PendingTasks = request.PendingTasks

	err = json.Unmarshal(request.ExternalDetails, &r.ExternalDetails)
	if err != nil {
		return nil, err
	}
	r.CreatedAt = regources.Time(request.CreatedAt)
	r.UpdatedAt = regources.Time(request.UpdatedAt)

	r.Details, err = PopulateDetails(request, request.RequestType, request.Details)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate reviewable request details")
	}

	return
}
