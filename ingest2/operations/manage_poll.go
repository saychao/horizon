package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type managePollOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *managePollOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	managePollOp := op.Body.MustManagePollOp()

	details := history2.OperationDetails{
		Type: xdr.OperationTypeManagePoll,
		ManagePoll: &history2.ManagePollDetails{
			Action: managePollOp.Data.Action,
			PollID: int64(managePollOp.PollId),
		},
	}

	switch managePollOp.Data.Action {
	case xdr.ManagePollActionClose:
		closeDetails := &history2.ClosePollData{
			Details:    internal.MarshalCustomDetails(managePollOp.Data.MustClosePollData().Details),
			PollResult: managePollOp.Data.MustClosePollData().Result,
		}
		details.ManagePoll.ClosePoll = closeDetails
	case xdr.ManagePollActionUpdateEndTime:
		updateEndTimeDetails := &history2.UpdatePollEndTimeData{
			EndTime: internal.TimeFromXdr(managePollOp.Data.UpdateTimeData.NewEndTime),
		}
		details.ManagePoll.UpdatePollEndTime = updateEndTimeDetails
	case xdr.ManagePollActionCancel:
	default:
		return history2.OperationDetails{}, errors.From(errors.New("unexpected manage poll action"),
			logan.F{
				"action": managePollOp.Data.Action.ShortString(),
			})
	}

	return details, nil
}
