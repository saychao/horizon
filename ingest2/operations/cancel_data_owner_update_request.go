package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type cancelDataOwnerUpdateRequestOpHandler struct {
	effectsProvider
}

func (h *cancelDataOwnerUpdateRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelDataOwnerUpdateRequestOp := op.Body.MustCancelDataOwnerUpdateRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelDataOwnerUpdateRequest,
		CancelDataOwnerUpdateRequest: &history2.CancelDataOwnerUpdateRequest{
			RequestID: uint64(cancelDataOwnerUpdateRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *cancelDataOwnerUpdateRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
