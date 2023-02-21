package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeOfferListBaseAssets - defines if the base assets should be included in the response
	IncludeTypeOfferListBaseAssets = "base_asset"
	// IncludeTypeOfferListQuoteAssets - defines if the quote assets should be included in the response
	IncludeTypeOfferListQuoteAssets = "quote_asset"

	// FilterTypeOfferListBaseBalance - defines if we need to filter the list by base balance
	FilterTypeOfferListBaseBalance = "base_balance"
	// FilterTypeOfferListQuoteBalance - defines if we need to filter the list by quote balance
	FilterTypeOfferListQuoteBalance = "quote_balance"
	// FilterTypeOfferListBaseAsset - defines if we need to filter the list by base asset
	FilterTypeOfferListBaseAsset = "base_asset"
	// FilterTypeOfferListQuoteAsset - defines if we need to filter the list by quote asset
	FilterTypeOfferListQuoteAsset = "quote_asset"
	// FilterTypeOfferListOwner - defines if we need to filter the list by owner
	FilterTypeOfferListOwner = "owner"
	// FilterTypeOfferListOrderBook - defines if we need to filter the list by order book
	FilterTypeOfferListOrderBook = "order_book"
	// FilterTypeOfferListIsBuy - defines if we need to filter the list by is buy
	FilterTypeOfferListIsBuy = "is_buy"
)

var includeTypeOfferListAll = map[string]struct{}{
	IncludeTypeOfferListBaseAssets:  {},
	IncludeTypeOfferListQuoteAssets: {},
}

var filterTypeOfferListAll = map[string]struct{}{
	FilterTypeOfferListBaseBalance:  {},
	FilterTypeOfferListQuoteBalance: {},
	FilterTypeOfferListBaseAsset:    {},
	FilterTypeOfferListQuoteAsset:   {},
	FilterTypeOfferListOwner:        {},
	FilterTypeOfferListOrderBook:    {},
	FilterTypeOfferListIsBuy:        {},
}

// GetOfferList represents params to be specified by user for getOfferList handler
type GetOfferList struct {
	*base
	Filters struct {
		BaseBalance  *string `filter:"base_balance"`
		QuoteBalance *string `filter:"quote_balance"`
		BaseAsset    *string `filter:"base_asset"`
		QuoteAsset   *string `filter:"quote_asset"`
		Owner        *string `filter:"owner"`
		OrderBook    *int64  `filter:"order_book"`
		IsBuy        *bool   `filter:"is_buy"`
	}
	PageParams pgdb.OffsetPageParams
	Includes   struct {
		BaseAsset   bool `include:"base_asset"`
		QuoteAssets bool `include:"quote_assets"`
	}
}

// NewGetOfferList - returns new instance of GetOfferList
func NewGetOfferList(r *http.Request) (*GetOfferList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeOfferListAll,
		supportedFilters:  filterTypeOfferListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetOfferList{
		base: b,
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
