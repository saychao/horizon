package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeCreateAmlAlertRequestsBalance = "request_details.balance"
)

var includeTypeCreateAmlAlertRequests = map[string]struct{}{
	IncludeTypeCreateAmlAlertRequestsBalance: {},
}

const (
	FilterTypeCreateAmlAlertRequestsBalance = "request_details.balance"
)

var filterTypeCreateAmlAlertRequests = map[string]struct{}{
	FilterTypeCreateAmlAlertRequestsBalance: {},
}

type GetCreateAmlAlertRequestsFilter struct {
	Balance *string `filter:"request_details.balance"`
}

type GetCreateAmlAlertRequests struct {
	GetRequestsBase
	Filters  GetCreateAmlAlertRequestsFilter
	Includes struct {
		RequestDetailsBalance bool `include:"request_details.balance"`
	}
}

func NewGetCreateAmlAlertRequests(r *http.Request) (request GetCreateAmlAlertRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCreateAmlAlertRequests,
		includeTypeCreateAmlAlertRequests,
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
