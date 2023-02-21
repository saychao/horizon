package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"

	"net/http"
)

const (
	// DeprecatedIncludeTypeOrderBookBaseAssets - defines if the base assets should be included in the response
	DeprecatedIncludeTypeOrderBookBaseAssets = "base_asset"
	// DeprecatedIncludeTypeOrderBookQuoteAssets = "quote_asset" - defines if the quote assets should be included in the response
	DeprecatedIncludeTypeOrderBookQuoteAssets = "quote_asset"

	// DeprecatedFilterTypeOrderBookBaseAsset - defines if we need to filter the list by base asset
	DeprecatedFilterTypeOrderBookBaseAsset = "base_asset"
	// DeprecatedFilterTypeOrderBookQuoteAsset - defines if we need to filter the list by quote asset
	DeprecatedFilterTypeOrderBookQuoteAsset = "quote_asset"
	// DeprecatedFilterTypeOrderBookIsBuy - defines if we need to filter the list by is buy
	DeprecatedFilterTypeOrderBookIsBuy = "is_buy"
)

var deprecatedIncludeTypeOrderBookAll = map[string]struct{}{
	DeprecatedIncludeTypeOrderBookBaseAssets:  {},
	DeprecatedIncludeTypeOrderBookQuoteAssets: {},
}

var deprecatedFilterTypeOrderBookAll = map[string]struct{}{
	DeprecatedFilterTypeOrderBookBaseAsset:  {},
	DeprecatedFilterTypeOrderBookQuoteAsset: {},
	DeprecatedFilterTypeOrderBookIsBuy:      {},
}

// DeprecatedGetOrderBook represents params to be specified by user for getOfferList handler
type DeprecatedGetOrderBook struct {
	*base
	ID      uint64
	Filters struct {
		BaseAsset  *string `filter:"base_asset"`
		QuoteAsset *string `filter:"quote_asset"`
		IsBuy      *bool   `filter:"is_buy"`
	}

	PageParams pgdb.OffsetPageParams
}

// NewDeprecatedGetOrderBook - returns new instance of DeprecatedGetOrderBook
func NewDeprecatedGetOrderBook(r *http.Request) (*DeprecatedGetOrderBook, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: deprecatedIncludeTypeOrderBookAll,
		supportedFilters:  deprecatedFilterTypeOrderBookAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("id")
	if err != nil {
		return nil, err
	}

	request := DeprecatedGetOrderBook{
		base: b,
		ID:   id,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = b.SetDefaultOffsetPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
