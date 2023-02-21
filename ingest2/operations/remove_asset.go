package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type removeAssetOpHandler struct {
	effectsProvider
}

// Details returns details about remove asset pair operation
func (h *removeAssetOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	removeAssetOp := op.Body.MustRemoveAssetOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeRemoveAsset,
		RemoveAsset: &history2.RemoveAssetDetails{
			Code: string(removeAssetOp.Code),
		},
	}, nil
}

// ParticipantsEffects - returns source of the operation
func (h *removeAssetOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
