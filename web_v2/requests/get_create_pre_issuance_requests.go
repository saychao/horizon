package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeCreatePreIssuanceRequestsAsset = "request_details.asset"
)

var filterTypeCreatePreIssuanceRequests = map[string]struct{}{
	FilterTypeCreatePreIssuanceRequestsAsset: {},
}

const (
	IncludeTypeCreatePreIssuanceRequestsAsset = "request_details.asset"
)

var includeTypeCreatePreIssuanceRequests = map[string]struct{}{
	IncludeTypeCreatePreIssuanceRequestsAsset: {},
}

type GetCreatePreIssuanceRequestsFilter struct {
	Asset *string `filter:"request_details.asset"`
}

type GetCreatePreIssuanceRequests struct {
	GetRequestsBase
	Filters  GetCreatePreIssuanceRequestsFilter
	Includes struct {
		RequestDetailsAsset bool `include:"request_details.asset"`
	}
}

func NewGetCreatePreIssuanceRequests(r *http.Request) (request GetCreatePreIssuanceRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreatePreIssuanceRequests,
		includeTypeCreatePreIssuanceRequests,
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
