package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageVoteOpHandler struct {
	effectsProvider
}

// Details returns details about manage vote operation
func (h *manageVoteOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageVoteOp := op.Body.MustManageVoteOp()

	details := history2.OperationDetails{
		Type: xdr.OperationTypeManageVote,
		ManageVote: &history2.ManageVoteDetails{
			Action: manageVoteOp.Data.Action,
		},
	}

	switch manageVoteOp.Data.Action {
	case xdr.ManageVoteActionCreate:
		details.ManageVote.PollID = int64(manageVoteOp.Data.MustCreateData().PollId)
		voteData := manageVoteOp.Data.MustCreateData().Data
		details.ManageVote.VoteData = &voteData
	case xdr.ManageVoteActionRemove:
		details.ManageVote.PollID = int64(manageVoteOp.Data.MustRemoveData().PollId)
	default:
		return history2.OperationDetails{}, errors.From(errors.New("unexpected manage vote action"),
			logan.F{
				"action": manageVoteOp.Data.Action.ShortString(),
			})
	}

	return details, nil
}
