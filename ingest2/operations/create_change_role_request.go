package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
)

type createChangeRoleRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create KYC request operation
func (h *createChangeRoleRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createKYCRequestOp := op.Body.MustCreateChangeRoleRequestOp()
	createKYCRequestRes := opRes.MustCreateChangeRoleRequestResult().MustSuccess()

	var allTasks *uint32
	if createKYCRequestOp.AllTasks != nil {
		allTasksInt := uint32(*createKYCRequestOp.AllTasks)
		allTasks = &allTasksInt
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateChangeRoleRequest,
		CreateChangeRoleRequest: &history2.CreateChangeRoleRequestDetails{
			DestinationAccount: createKYCRequestOp.DestinationAccount.Address(),
			AccountRoleToSet:   uint64(createKYCRequestOp.AccountRoleToSet),
			CreatorDetails:     internal.MarshalCustomDetails(createKYCRequestOp.CreatorDetails),
			AllTasks:           allTasks,
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(createKYCRequestRes.RequestId),
				IsFulfilled: createKYCRequestRes.Fulfilled,
			},
		},
	}, nil
}

// ParticipantsEffects returns source participant effect and effect for account for which KYC is updated
func (h *createChangeRoleRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	createKYCRequestOp := opBody.MustCreateChangeRoleRequestOp()

	updatedAccount := h.Participant(createKYCRequestOp.DestinationAccount)

	source := h.Participant(sourceAccountID)
	if updatedAccount.AccountID == source.AccountID {
		return []history2.ParticipantEffect{source}, nil
	}

	return []history2.ParticipantEffect{source, updatedAccount}, nil
}
