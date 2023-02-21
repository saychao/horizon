package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type cancelDataUpdateRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *cancelDataUpdateRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelDataUpdateRequestOp := op.Body.MustCancelDataUpdateRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelDataUpdateRequest,
		CancelDataUpdateRequest: &history2.CancelDataUpdateRequest{
			RequestID: uint64(cancelDataUpdateRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *cancelDataUpdateRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
