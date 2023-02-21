package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeCreateWithdrawRequestsBalance = "request_details.balance"
	IncludeTypeCreateWithdrawRequestsAsset   = "request_details.asset"
)

var includeTypeCreateWithdrawRequests = map[string]struct{}{
	IncludeTypeCreateWithdrawRequestsBalance: {},
	IncludeTypeCreateWithdrawRequestsAsset:   {},
}

const (
	FilterTypeCreateWithdrawRequestsBalance = "request_details.balance"
	FilterTypeCreateWithdrawRequestsAsset   = "request_details.asset"
)

var filterTypeCreateWithdrawRequests = map[string]struct{}{
	FilterTypeCreateWithdrawRequestsBalance: {},
	FilterTypeCreateWithdrawRequestsAsset:   {},
}

type GetCreateWithdrawRequestsFilter struct {
	Balance *string  `filter:"request_details.balance"`
	Asset   []string `filter:"request_details.asset"`
}

type GetCreateWithdrawRequests struct {
	GetRequestsBase
	Filters  GetCreateWithdrawRequestsFilter
	Includes struct {
		Balance bool `include:"request_details.balance"`
		Asset   bool `include:"request_details.asset"`
	}
}

func NewGetCreateWithdrawRequests(r *http.Request) (request GetCreateWithdrawRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateWithdrawRequests,
		includeTypeCreateWithdrawRequests,
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
