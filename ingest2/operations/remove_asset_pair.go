package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type removeAssetPairOpHandler struct {
	effectsProvider
}

// Details returns details about remove asset pair operation
func (h *removeAssetPairOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	removeAssetPairOp := op.Body.MustRemoveAssetPairOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeRemoveAssetPair,
		RemoveAssetPair: &history2.RemoveAssetPairDetails{
			Base:  string(removeAssetPairOp.Base),
			Quote: string(removeAssetPairOp.Quote),
		},
	}, nil
}

// ParticipantsEffects - returns source of the operation
func (h *removeAssetPairOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
