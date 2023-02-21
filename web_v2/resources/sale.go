package resources

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/amount"
	regources "gitlab.com/tokend/regources/generated"
)

// NewSaleKey - creates new Key for asset
func NewSaleKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, regources.SALES)
}

// NewSale - creates new instance of Sale
func NewSale(record history2.Sale) regources.Sale {
	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(record.QuoteAssets.QuoteAssets)),
	}

	for _, quoteAsset := range record.QuoteAssets.QuoteAssets {
		quoteAssets.Data = append(quoteAssets.Data, newQuoteAssetKey(quoteAsset.Asset))
	}

	return regources.Sale{
		Key: NewSaleKey(int64(record.ID)),
		Attributes: regources.SaleAttributes{
			AccessDefinitionType: record.AccessDefinitionType,
			BaseHardCap:          record.BaseHardCap,
			StartTime:            record.StartTime,
			EndTime:              record.EndTime,
			SaleType:             record.SaleType,
			SaleState:            record.State,
			Details:              record.Details,
		},
		Relationships: regources.SaleRelationships{
			Owner:             NewAccountKey(record.OwnerAddress).AsRelation(),
			BaseAsset:         NewAssetKey(record.BaseAsset).AsRelation(),
			QuoteAssets:       quoteAssets,
			DefaultQuoteAsset: newQuoteAssetKey(record.DefaultQuoteAsset).AsRelation(),
		},
	}
}

// NewSaleDefaultQuoteAsset - extracts the default quote asset details from the sale
func NewSaleDefaultQuoteAsset(saleRecord history2.Sale) regources.SaleQuoteAsset {
	var price regources.Amount = amount.One

	for _, quoteAsset := range saleRecord.QuoteAssets.QuoteAssets {
		if quoteAsset.Asset == saleRecord.DefaultQuoteAsset {
			price = quoteAsset.Price
		}
	}

	return regources.SaleQuoteAsset{
		Key: *NewSaleQuoteAssetKey(saleRecord.DefaultQuoteAsset, saleRecord.ID).GetKeyP(),
		Attributes: regources.SaleQuoteAssetAttributes{
			Price:      price,
			CurrentCap: saleRecord.CurrentCap,
			HardCap:    saleRecord.HardCap,
			SoftCap:    saleRecord.SoftCap,
		},
		Relationships: regources.SaleQuoteAssetRelationships{
			Asset: NewAssetKey(saleRecord.DefaultQuoteAsset).AsRelation(),
		},
	}
}
