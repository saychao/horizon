package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type manageRemoveDataOpHandler struct {
	effectsProvider
}

func (h *manageRemoveDataOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceID)}, nil
}

func (h *manageRemoveDataOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	removeDataOp := op.Body.MustRemoveDataOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeRemoveData,
		RemoveData: &history2.RemoveDataDetails{
			ID: uint64(removeDataOp.DataId),
		},
	}, nil
}
