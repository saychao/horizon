package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type stubOpHandler struct {
	effectsProvider
}

// Details - used as temporary solution for not handled operations
func (h *stubOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr) (history2.OperationDetails, error) {
	return history2.OperationDetails{
		Type: op.Body.Type,
	}, nil
}

// ParticipantsEffects - used as temroary solution for not handled operations. Returns only source participant
func (h *stubOpHandler) ParticipantsEffects(opBody xdr.OperationBody, opRes xdr.OperationResultTr,
	sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
