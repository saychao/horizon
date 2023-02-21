package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeCreateAtomicSwapAskRequestsBalance     = "request_details.base_balance"
	IncludeTypeCreateAtomicSwapAskRequestsQuoteAssets = "request_details.quote_assets"
)

var includeTypeCreateAtomicSwapAskRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapAskRequestsBalance:     {},
	IncludeTypeCreateAtomicSwapAskRequestsQuoteAssets: {},
}

const (
	FilterTypeCreateAtomicSwapAskRequestsBalance = "request_details.base_balance"
)

var filterTypeCreateAtomicSwapAskRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapAskRequestsBalance: {},
}

type GetCreateAtomicSwapAskRequestsFilter struct {
	BaseBalance *string `filter:"request_details.base_balance"`
}

type GetCreateAtomicSwapAskRequests struct {
	GetRequestsBase
	Filters  GetCreateAtomicSwapAskRequestsFilter
	Includes struct {
		RequestDetailsBaseBalance bool `include:"request_details.base_balance"`
		RequestDetailsQuoteAssets bool `include:"request_details.quote_assets"`
	}
}

func NewGetCreateAtomicSwapAskRequests(r *http.Request) (request GetCreateAtomicSwapAskRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateAtomicSwapAskRequests,
		includeTypeCreateAtomicSwapAskRequests,
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
