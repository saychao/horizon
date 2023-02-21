package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeCreateIssuanceRequestsAsset    = "request_details.asset"
	FilterTypeCreateIssuanceRequestsReceiver = "request_details.receiver"
	FilterTypeCreateIssuanceRequestsWho      = "request_details.who"
)

var filterTypeCreateIssuanceRequests = map[string]struct{}{
	FilterTypeCreateIssuanceRequestsAsset:    {},
	FilterTypeCreateIssuanceRequestsReceiver: {},
	FilterTypeCreateIssuanceRequestsWho:      {},
}

const (
	IncludeTypeCreateIssuanceRequestsAsset = "request_details.asset"
)

var includeTypeCreateIssuanceRequests = map[string]struct{}{
	IncludeTypeCreateIssuanceRequestsAsset: {},
}

type GetCreateIssuanceRequestsFilter struct {
	Asset    []string `filter:"request_details.asset"`
	Receiver *string  `filter:"request_details.receiver"`
	Who      *string  `filter:"request_details.who"`
}

type GetCreateIssuanceRequests struct {
	GetRequestsBase
	Filters  GetCreateIssuanceRequestsFilter
	Includes struct {
		RequestDetailsAsset bool `include:"request_details.asset"`
	}
}

func NewGetCreateIssuanceRequests(r *http.Request) (request GetCreateIssuanceRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateIssuanceRequests,
		includeTypeCreateIssuanceRequests,
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
