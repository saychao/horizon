package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset = "request_details.quote_asset"
)

var includeTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
}

const (
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset = "request_details.quote_asset"
	FilterTypeCreateAtomicSwapBidRequestsAskID      = "request_details.ask_id"
	FilterTypeCreateAtomicSwapBidRequestsAskOwner   = "request_details.ask_owner"
)

var filterTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
	FilterTypeCreateAtomicSwapBidRequestsAskID:      {},
	FilterTypeCreateAtomicSwapBidRequestsAskOwner:   {},
}

type GetCreateAtomicSwapBidRequestsFilter struct {
	QuoteAsset *string `filter:"request_details.quote_asset"`
	AskID      *uint64 `filter:"request_details.ask_id"`
	AskOwner   *string `filter:"request_details.ask_owner"`
}

type GetCreateAtomicSwapBidRequests struct {
	GetRequestsBase
	Filters  GetCreateAtomicSwapBidRequestsFilter
	Includes struct {
		RequestDetailsQuoteAssets bool `include:"request_details.quote_assets"`
	}
}

func NewGetCreateAtomicSwapBidRequests(r *http.Request) (request GetCreateAtomicSwapBidRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateAtomicSwapBidRequests,
		includeTypeCreateAtomicSwapBidRequests,
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
