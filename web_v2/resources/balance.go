package resources

import (
	core "github.com/saychao/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewBalance - creates new instance of balance using core balance
func NewBalance(record *core.Balance) *regources.Balance {
	return &regources.Balance{
		Key: NewBalanceKey(record.BalanceAddress),
	}
}

// NewBalanceKey - creates new instance of balance key
func NewBalanceKey(balanceAddress string) regources.Key {
	return regources.Key{
		Type: regources.BALANCES,
		ID:   balanceAddress,
	}
}

// NewBalanceState - returns new balance state
func NewBalanceState(record *core.Balance) *regources.BalanceState {
	return &regources.BalanceState{
		Key: regources.Key{
			ID:   record.BalanceAddress,
			Type: regources.BALANCES_STATE,
		},
		Attributes: &regources.BalanceStateAttributes{
			Locked:    regources.Amount(record.Locked),
			Available: regources.Amount(record.Amount),
		},
	}
}

// NewBalanceStateKey - creates new balance state key using balance address
func NewBalanceStateKey(balanceAddress string) regources.Key {
	return regources.Key{
		ID:   balanceAddress,
		Type: regources.BALANCES_STATE,
	}
}
