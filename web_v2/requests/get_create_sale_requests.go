package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeCreateSaleRequestsBaseAsset         = "request_details.base_asset"
	FilterTypeCreateSaleRequestsDefaultQuoteAsset = "request_details.default_quote_asset"
)

const (
	IncludeTypeCreateSaleRequestsBaseAsset         = "request_details.base_asset"
	IncludeTypeCreateSaleRequestsDefaultQuoteAsset = "request_details.default_quote_asset"
	IncludeTypeCreateSaleRequestsQuoteAssets       = "request_details.quote_assets"
)

var filterTypeCreateSaleRequests = map[string]struct{}{
	FilterTypeCreateSaleRequestsBaseAsset:         {},
	FilterTypeCreateSaleRequestsDefaultQuoteAsset: {},
}

var includeTypeCreateSaleRequests = map[string]struct{}{
	IncludeTypeCreateSaleRequestsBaseAsset:         {},
	IncludeTypeCreateSaleRequestsDefaultQuoteAsset: {},
	IncludeTypeCreateSaleRequestsQuoteAssets:       {},
}

type GetCreateSaleRequestsFilter struct {
	BaseAsset         *string `filter:"request_details.base_asset"`
	DefaultQuoteAsset *string `filter:"request_details.default_quote_asset"`
}

type GetCreateSaleRequests struct {
	GetRequestsBase
	Filters  GetCreateSaleRequestsFilter
	Includes struct {
		BaseAsset         bool `include:"request_details.base_asset"`
		DefaultQuoteAsset bool `include:"request_details.default_quote_asset"`
		QuoteAssets       bool `include:"request_details.quote_assets"`
	}
}

func NewGetCreateSaleRequests(r *http.Request) (request GetCreateSaleRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateSaleRequests,
		includeTypeCreateSaleRequests,
	)
	if err != nil {
		return request, err
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	err = PopulateRequest(&request.GetRequestsBase)
	if err != nil {
		return request, err
	}

	return request, nil
}
