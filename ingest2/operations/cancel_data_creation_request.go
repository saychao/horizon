package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type cancelDataCreationRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *cancelDataCreationRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelDataCreationRequestOp := op.Body.MustCancelDataCreationRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelDataCreationRequest,
		CancelDataCreationRequest: &history2.CancelDataCreationRequest{
			RequestID: uint64(cancelDataCreationRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *cancelDataCreationRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
