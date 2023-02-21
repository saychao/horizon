package requests

import "net/http"

const (
	// IncludeTypeSaleBaseAsset - defines if the base asset should be included in the response
	IncludeTypeSaleBaseAsset = "base_asset"
	// IncludeTypeSaleQuoteAssets - defines if the base asset should be included in the response
	IncludeTypeSaleQuoteAssets = "quote_assets"
	// IncludeTypeSaleDefaultQuoteAsset - defines if the default quote asset should be included in the response
	IncludeTypeSaleDefaultQuoteAsset = "default_quote_asset"
)

var includeTypeSaleAll = map[string]struct{}{
	IncludeTypeSaleBaseAsset:         {},
	IncludeTypeSaleQuoteAssets:       {},
	IncludeTypeSaleDefaultQuoteAsset: {},
}

// GetSale represents params to be specified by user for getSale handler
type GetSale struct {
	*base
	ID       uint64
	Includes struct {
		BaseAsset         bool `include:"base_asset"`
		QuoteAssets       bool `include:"quote_assets"`
		DefaultQuoteAsset bool `include:"default_quote_asset"`
	}
}

// NewGetSale returns new instance of the GetSale request
func NewGetSale(r *http.Request) (*GetSale, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("id")
	if err != nil {
		return nil, err
	}

	return &GetSale{
		base: b,
		ID:   id,
	}, nil
}
