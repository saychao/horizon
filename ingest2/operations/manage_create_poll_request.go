package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageCreatePollRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *manageCreatePollRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageCreatePollRequestOp := op.Body.MustManageCreatePollRequestOp()
	manageCreatePollRequestOpRes := opRes.MustManageCreatePollRequestResult().MustSuccess()

	details := history2.OperationDetails{
		Type: xdr.OperationTypeManageCreatePollRequest,
		ManageCreatePollRequest: &history2.ManageCreatePollRequestDetails{
			Action: manageCreatePollRequestOp.Data.Action,
		},
	}

	switch manageCreatePollRequestOp.Data.Action {
	case xdr.ManageCreatePollRequestActionCreate:
		var allTasks *uint32
		rawAllTasks := manageCreatePollRequestOp.Data.CreateData.AllTasks
		if rawAllTasks != nil {
			allTasksInt := uint32(*rawAllTasks)
			allTasks = &allTasksInt
		}

		createPollRequest := manageCreatePollRequestOp.Data.MustCreateData().Request
		createDetails := &history2.CreatePollRequestDetails{
			AllTasks:                 allTasks,
			PollType:                 createPollRequest.Data.Type,
			VoteConfirmationRequired: createPollRequest.VoteConfirmationRequired,
			ResultProviderID:         createPollRequest.ResultProviderId.Address(),
			PermissionType:           uint32(createPollRequest.PermissionType),
			NumberOfChoices:          uint32(createPollRequest.NumberOfChoices),
			StartTime:                internal.TimeFromXdr(createPollRequest.StartTime),
			EndTime:                  internal.TimeFromXdr(createPollRequest.EndTime),
			CreatorDetails:           internal.MarshalCustomDetails(createPollRequest.CreatorDetails),
			PollData:                 createPollRequest.Data,
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(manageCreatePollRequestOpRes.Details.MustResponse().RequestId),
				IsFulfilled: manageCreatePollRequestOpRes.Details.MustResponse().Fulfilled,
			},
		}

		details.ManageCreatePollRequest.CreateDetails = createDetails

	case xdr.ManageCreatePollRequestActionCancel:
		cancelDetails := &history2.CancelCreatePollRequestDetails{
			RequestID: int64(manageCreatePollRequestOp.Data.MustCancelData().RequestId),
		}
		details.ManageCreatePollRequest.CancelDetails = cancelDetails
	default:
		return history2.OperationDetails{}, errors.From(errors.New("unexpected manage create poll request action"),
			logan.F{
				"action": manageCreatePollRequestOp.Data.Action.ShortString(),
			})
	}

	return details, nil
}

// ParticipantsEffects - returns source of the operation and account for which balance was created if they differ
func (h *manageCreatePollRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageCreatePollRequestOp := opBody.MustManageCreatePollRequestOp()

	var participants []history2.ParticipantEffect
	source := h.Participant(sourceAccountID)

	if manageCreatePollRequestOp.Data.Action != xdr.ManageCreatePollRequestActionCreate {
		return []history2.ParticipantEffect{source}, nil
	}

	resultProvider := h.Participant(manageCreatePollRequestOp.Data.MustCreateData().Request.ResultProviderId)
	if source.AccountID != resultProvider.AccountID {
		participants = []history2.ParticipantEffect{resultProvider}
	}

	return append(participants, source), nil
}
