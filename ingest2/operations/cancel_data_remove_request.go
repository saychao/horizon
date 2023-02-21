package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type cancelDataRemoveRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *cancelDataRemoveRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelDataRemoveRequestOp := op.Body.MustCancelDataRemoveRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelDataRemoveRequest,
		CancelDataRemoveRequest: &history2.CancelDataRemoveRequest{
			RequestID: uint64(cancelDataRemoveRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *cancelDataRemoveRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
