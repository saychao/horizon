package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeUpdateAssetRequestsAsset = "request_details.asset"
)

var filterTypeUpdateAssetRequests = map[string]struct{}{
	FilterTypeUpdateAssetRequestsAsset: {},
}

const (
	IncludeTypeUpdateAssetRequestsAsset = "request_details.asset"
)

var includeTypeUpdateAssetRequests = map[string]struct{}{
	IncludeTypeUpdateAssetRequestsAsset: {},
}

type GetUpdateAssetRequestsFilter struct {
	Asset *string `filter:"request_details.asset"`
}

type GetUpdateAssetRequests struct {
	GetRequestsBase
	Filters  GetUpdateAssetRequestsFilter
	Includes struct {
		RequestDetailsAsset bool `include:"request_details.asset"`
	}
}

func NewGetUpdateAssetRequests(r *http.Request) (request GetUpdateAssetRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeUpdateAssetRequests,
		includeTypeUpdateAssetRequests,
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
