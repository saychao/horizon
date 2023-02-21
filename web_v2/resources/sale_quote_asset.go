package resources

import (
	"fmt"
	"github.com/saychao/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewSaleQuoteAssetKey - returns new instance of SaleQuoteAssetKey
func NewSaleQuoteAssetKey(assetCode string, saleID uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprintf("%s:%d", assetCode, saleID),
		Type: regources.SALE_QUOTE_ASSETS,
	}
}

// NewSaleQuoteAsset - returns new instance of SaleQuoteAsset
func NewSaleQuoteAsset(qa history2.SaleQuoteAsset, saleID uint64) regources.SaleQuoteAsset {
	return regources.SaleQuoteAsset{
		Key: NewSaleQuoteAssetKey(qa.Asset, saleID),
		Attributes: regources.SaleQuoteAssetAttributes{
			Price:           qa.Price,
			CurrentCap:      qa.CurrentCap,
			TotalCurrentCap: qa.TotalCurrentCap,
			HardCap:         qa.HardCap,
		},
		Relationships: regources.SaleQuoteAssetRelationships{
			Asset: NewAssetKey(qa.Asset).AsRelation(),
		},
	}
}
