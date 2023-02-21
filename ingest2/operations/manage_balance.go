package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type manageBalanceOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *manageBalanceOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageBalanceOp := op.Body.MustManageBalanceOp()
	manageBalanceRes := opRes.MustManageBalanceResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageBalance,
		ManageBalance: &history2.ManageBalanceDetails{
			DestinationAccount: manageBalanceOp.Destination.Address(),
			Action:             manageBalanceOp.Action,
			Asset:              string(manageBalanceOp.Asset),
			BalanceAddress:     manageBalanceRes.BalanceId.AsString(),
		},
	}, nil
}

// ParticipantsEffects - returns source of the operation and account for which balance was created if they differ
func (h *manageBalanceOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageBalanceOp := opBody.MustManageBalanceOp()

	destinationAccount := h.Participant(manageBalanceOp.Destination)
	source := h.Participant(sourceAccountID)

	var participants []history2.ParticipantEffect

	if source.AccountID != destinationAccount.AccountID {
		participants = []history2.ParticipantEffect{destinationAccount}
	}

	return append(participants, source), nil
}
