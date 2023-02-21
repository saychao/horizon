package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type createAccountOpHandler struct {
	effectsProvider
}

// Details returns details about create account operation
func (h *createAccountOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createAccOp := op.Body.MustCreateAccountOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAccount,
		CreateAccount: &history2.CreateAccountDetails{
			AccountAddress: createAccOp.Destination.Address(),
			AccountRole:    uint64(createAccOp.RoleId),
		},
	}, nil
}

// ParticipantsEffects returns counterparties without effects
func (h *createAccountOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, changes []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	source := h.Participant(sourceAccountID)
	participants := []history2.ParticipantEffect{source}

	createAccountOp := opBody.MustCreateAccountOp()

	participants = append(participants, history2.ParticipantEffect{
		AccountID: h.MustAccountID(createAccountOp.Destination),
	})

	referrerEffect, err := h.referrerParticipantEffect(createAccountOp, changes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get referrer participant effect")
	}

	if referrerEffect != nil {
		participants = append(participants, *referrerEffect)
	}

	return participants, nil
}

// referrerParticipantEffect - provides referrer participant effect for create account op
// handles case when due to some reasons in createAccountOp we have received ID of non existing account
func (h *createAccountOpHandler) referrerParticipantEffect(op xdr.CreateAccountOp,
	changes []xdr.LedgerEntryChange) (*history2.ParticipantEffect, error) {
	if op.Referrer == nil {
		return nil, nil
	}

	account := mustFindNewAccountEntry(op.Destination, changes)
	if account.Referrer == nil {
		return nil, nil
	}

	return &history2.ParticipantEffect{
		AccountID: h.MustAccountID(*account.Referrer),
	}, nil
}

func mustFindNewAccountEntry(expectedAccountID xdr.AccountId, changes []xdr.LedgerEntryChange) xdr.AccountEntry {
	expectedAccountAddr := expectedAccountID.Address()
	for _, change := range changes {
		if change.Type != xdr.LedgerEntryChangeTypeCreated {
			continue
		}

		created := change.MustCreated()
		if created.Data.Type != xdr.LedgerEntryTypeAccount {
			continue
		}

		account := created.Data.MustAccount()
		if expectedAccountAddr == account.AccountId.Address() {
			return account
		}
	}

	panic(errors.From(errors.New("failed to find created account in createAccountOp changes"), logan.F{
		"expected_account_addr": expectedAccountAddr,
	}))
}
