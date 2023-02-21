package resources

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/regources/generated"
)

// NewMatchKey - creates new Key for OrderBook
func NewMatchKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, regources.MATCHES)
}

// NewMatch - returns new instance of `Match` resource
func NewMatch(record history2.Match) regources.Match {
	return regources.Match{
		Key: NewMatchKey(record.ID),
		Attributes: regources.MatchAttributes{
			BaseAmount:  regources.Amount(record.BaseAmount),
			Price:       regources.Amount(record.Price),
			QuoteAmount: regources.Amount(record.QuoteAmount),
			CreatedAt:   record.CreatedAt,
		},
		Relationships: regources.MatchRelationships{
			BaseAsset:  *NewAssetKey(record.BaseAsset).AsRelation(),
			QuoteAsset: *NewAssetKey(record.QuoteAsset).AsRelation(),
		},
	}
}
