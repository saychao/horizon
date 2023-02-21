package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type updateDataOwnerOpHandler struct {
	effectsProvider
}

func (h *updateDataOwnerOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{
		h.Participant(sourceID),
		h.Participant(opBody.UpdateDataOwnerOp.NewOwner),
	}, nil
}

func (h *updateDataOwnerOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	updateDataOwnerOp := op.Body.MustUpdateDataOwnerOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeUpdateData,
		UpdateDataOwner: &history2.UpdateDataOwnerDetails{
			NewOwner: updateDataOwnerOp.NewOwner,
			ID:       uint64(updateDataOwnerOp.DataId),
		},
	}, nil
}
