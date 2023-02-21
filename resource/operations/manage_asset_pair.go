package operations

type ManageAssetPair struct {
	Base
	BaseAsset               string `json:"base_asset"`
	QuoteAsset              string `json:"quote_asset"`
	PhysicalPrice           string `json:"physical_price"`
	PhysicalPriceCorrection string `json:"physical_price_correction"`
	MaxPriceStep            string `json:"max_price_step"`
}
