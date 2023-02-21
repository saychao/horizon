package resources

import regources "gitlab.com/tokend/regources/generated"

// newQuoteAssetKey - creates new key for resource type `quote-assets`
func newQuoteAssetKey(quoteAssetCode string) regources.Key {
	return regources.Key{
		ID:   quoteAssetCode,
		Type: regources.QUOTE_ASSETS,
	}
}

// newQuoteAssetKey - creates new key for resource type `sale-quote-assets`
func newSaleQuoteAssetKey(quoteAssetCode string) regources.Key {
	return regources.Key{
		ID:   quoteAssetCode,
		Type: regources.SALE_QUOTE_ASSETS,
	}
}
