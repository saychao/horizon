package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeDataCreationRequestsSecurityType = "request_details.security_type"
)

var filterTypeDataCreationRequests = map[string]struct{}{
	FilterTypeDataCreationRequestsSecurityType: {},
}

type GetDataCreationRequestsFilter struct {
	SecurityType *uint32 `filter:"request_details.security_type"`
}

type GetDataCreationRequests struct {
	GetRequestsBase
	Filters GetDataCreationRequestsFilter
}

func NewGetDataCreationRequests(r *http.Request) (request GetDataCreationRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeDataCreationRequests,
		map[string]struct{}{},
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
