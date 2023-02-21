package resources

import (
	"github.com/saychao/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewRequestKey - creates new instance of request key
func NewRequestKey(requestId int64) regources.Key {
	return regources.NewKeyInt64(requestId, regources.REQUESTS)
}

// NewRequest - creates new instance of reviewable request
func NewRequest(record history2.ReviewableRequest) regources.ReviewableRequest {
	return regources.ReviewableRequest{
		Key: NewRequestKey(record.ID),
		Attributes: regources.ReviewableRequestAttributes{
			Reference:       record.Reference,
			RejectReason:    record.RejectReason,
			Hash:            record.Hash,
			AllTasks:        record.AllTasks,
			PendingTasks:    record.PendingTasks,
			ExternalDetails: record.ExternalDetails,
			CreatedAt:       record.CreatedAt,
			UpdatedAt:       record.UpdatedAt,
			XdrType:         record.RequestType,

			// TODO shouldn't those look like `state: {str: "", int: 2}`?
			//  or just drop int part?
			State:  record.RequestState.String(),
			StateI: int32(record.RequestState),
		},
		Relationships: regources.ReviewableRequestRelationships{
			Requestor: NewAccountKey(record.Requestor).AsRelation(),
			Reviewer:  NewAccountKey(record.Reviewer).AsRelation(),
		},
	}
}
