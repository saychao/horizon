package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeCreateAssetRequestsAsset = "request_details.asset"
)

var filterTypeCreateAssetRequests = map[string]struct{}{
	FilterTypeCreateAssetRequestsAsset: {},
}

const (
	IncludeTypeCreateAssetRequestsAsset = "request_details.asset"
)

var includeTypeCreateAssetRequests = map[string]struct{}{
	IncludeTypeCreateAssetRequestsAsset: {},
}

type GetCreateAssetRequestsFilter struct {
	Asset *string `filter:"request_details.asset"`
}

type GetCreateAssetRequests struct {
	GetRequestsBase
	Filters  GetCreateAssetRequestsFilter
	Includes struct {
		RequestDetailsAsset bool `include:"request_details.asset"`
	}
}

func NewGetCreateAssetRequests(r *http.Request) (request GetCreateAssetRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateAssetRequests,
		includeTypeCreateAssetRequests,
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
