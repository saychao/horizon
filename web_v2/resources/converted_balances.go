package resources

import (
	"github.com/saychao/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewConvertedBalancesCollectionKey - creates the key of ConvertedBalancesCollection resource
func NewConvertedBalancesCollectionKey(assetCode string) regources.Key {
	return regources.Key{
		ID:   assetCode,
		Type: regources.CONVERTED_BALANCES_COLLECTIONS,
	}
}

// NewConvertedBalanceCollection - creates new instance of ConvertedBalancesCollection resource
func NewConvertedBalanceCollection(assetCode string) regources.ConvertedBalancesCollection {
	return regources.ConvertedBalancesCollection{
		Key: NewConvertedBalancesCollectionKey(assetCode),
		Relationships: regources.ConvertedBalancesCollectionRelationships{
			Asset: *NewAssetKey(assetCode).AsRelation(),
		},
	}
}

// NewConvertedBalanceStateKey - creates new instance of ConvertedBalanceState resource
func NewConvertedBalanceStateKey(balanceAddress string) regources.Key {
	return regources.Key{
		ID:   balanceAddress,
		Type: regources.CONVERTED_BALANCE_STATES,
	}
}

// NewConvertedBalanceState - creates new instance of ConvertedBalanceState resource
func NewConvertedBalanceState(
	balance core2.Balance,
	convertedAvailable regources.Amount,
	convertedLocked regources.Amount,
	price regources.Amount,
	isConverted bool,
) regources.ConvertedBalanceState {

	return regources.ConvertedBalanceState{
		Key: NewConvertedBalanceStateKey(balance.BalanceAddress),
		Attributes: regources.ConvertedBalanceStateAttributes{
			InitialAmounts: regources.BalanceStateAttributeAmounts{
				Available: regources.Amount(balance.Amount),
				Locked:    regources.Amount(balance.Locked),
			},
			ConvertedAmounts: regources.BalanceStateAttributeAmounts{
				Available: convertedAvailable,
				Locked:    convertedLocked,
			},
			IsConverted: isConverted,
			Price:       price,
		},
		Relationships: regources.ConvertedBalanceStateRelationships{
			Balance: NewBalanceKey(balance.BalanceAddress).AsRelation(),
		},
	}
}
