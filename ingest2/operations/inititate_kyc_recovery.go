package operations

import (
	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

type initiateKycRecoveryOpHandler struct {
	effectsProvider
}

// Details returns details about create KYC request operation
func (h *initiateKycRecoveryOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	initiateKycRecoveryOp := op.Body.MustInitiateKycRecoveryOp()

	publicKey := xdr.AccountId(initiateKycRecoveryOp.Signer)

	return history2.OperationDetails{
		Type: xdr.OperationTypeInitiateKycRecovery,
		InitiateKYCRecovery: &history2.InitiateKYCRecoveryDetails{
			Account: initiateKycRecoveryOp.Account.Address(),
			Signer:  publicKey.Address(),
		},
	}, nil
}

//ParticipantsEffects returns source participant effect and effect for account for which KYC is updated
func (h *initiateKycRecoveryOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	initiateKycRecoveryOp := opBody.MustInitiateKycRecoveryOp()

	updatedAccount := h.Participant(initiateKycRecoveryOp.Account)

	source := h.Participant(sourceAccountID)
	if updatedAccount.AccountID == source.AccountID {
		return []history2.ParticipantEffect{source}, nil
	}

	return []history2.ParticipantEffect{source, updatedAccount}, nil
}
