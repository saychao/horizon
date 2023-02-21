package resources

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewSwapKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.SWAPS)
}

func NewSwap(record history2.Swap) regources.Swap {
	return regources.Swap{
		Key: NewSwapKey(record.ID),
		Attributes: regources.SwapAttributes{
			Amount:    regources.Amount(record.Amount),
			CreatedAt: record.CreatedAt,
			DestinationFee: regources.Fee{
				Fixed:             regources.Amount(record.DestinationFixedFee),
				CalculatedPercent: regources.Amount(record.DestinationPercentFee),
			},
			Details:    record.Details,
			LockTime:   record.LockTime,
			Secret:     record.Secret,
			SecretHash: record.SecretHash,
			SourceFee: regources.Fee{
				Fixed:             regources.Amount(record.SourceFixedFee),
				CalculatedPercent: regources.Amount(record.SourcePercentFee),
			},
			State: record.State,
		},
		Relationships: regources.SwapRelationships{
			Asset:              NewAssetKey(record.Asset).AsRelation(),
			Destination:        NewAccountKey(record.DestinationAccount).AsRelation(),
			DestinationBalance: NewBalanceKey(record.DestinationBalance).AsRelation(),
			Source:             NewAccountKey(record.SourceAccount).AsRelation(),
			SourceBalance:      NewBalanceKey(record.SourceBalance).AsRelation(),
		},
	}
}
