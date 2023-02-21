package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeDeferredPaymentCreateRequestsDestination = "request_details.destination"
)

var filterTypeDeferredPaymentCreateRequests = map[string]struct{}{
	FilterTypeDeferredPaymentCreateRequestsDestination: {},
}

type GetDeferredPaymentCreateRequestsFilter struct {
	Destination *string `filter:"request_details.destination"`
}

type GetDeferredPaymentCreateRequests struct {
	GetRequestsBase
	Filters GetDeferredPaymentCreateRequestsFilter
}

func NewGetDeferredPaymentCreateRequests(r *http.Request) (request GetDeferredPaymentCreateRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeDeferredPaymentCreateRequests,
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
