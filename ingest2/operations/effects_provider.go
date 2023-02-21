package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

// effectsProvider - provides default implementation of ParticipantsEffects and helper methods to safely
// build ParticipantEffects
type effectsProvider struct {
	IDProvider
	balanceProvider
}

// ParticipantsEffects returns only source without effects
func (h *effectsProvider) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}

// Participant - returns participant without any effects
func (h *effectsProvider) Participant(sourceAccountID xdr.AccountId) history2.ParticipantEffect {
	return history2.ParticipantEffect{
		AccountID: h.MustAccountID(sourceAccountID),
	}
}

// BalanceEffect - creates participant effect for specified balance
func (h *effectsProvider) BalanceEffect(balanceID xdr.BalanceId,
	effect *history2.Effect) history2.ParticipantEffect {
	balance := h.MustBalance(balanceID)
	return history2.ParticipantEffect{
		AccountID: balance.AccountID,
		BalanceID: &balance.ID,
		AssetCode: &balance.AssetCode,
		Effect:    effect,
	}
}

// BalanceEffectWithAccount - creates participant effects slice for specified balance and account
// if account is owner of specified balance - returns only one participant
// returns account participant and balance participant otherwise
func (h *effectsProvider) BalanceEffectWithAccount(accountID xdr.AccountId, balanceID xdr.BalanceId,
	effect *history2.Effect) []history2.ParticipantEffect {
	account := h.Participant(accountID)
	balance := h.BalanceEffect(balanceID, effect)
	if account.AccountID == balance.AccountID {
		return []history2.ParticipantEffect{balance}
	}

	return []history2.ParticipantEffect{account, balance}
}
