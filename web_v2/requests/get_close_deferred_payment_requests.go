package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

var filterTypeCloseDeferredPaymentRequests = map[string]struct{}{}

type GetCloseDeferredPaymentRequestsFilter struct {
}

type GetCloseDeferredPaymentRequests struct {
	GetRequestsBase
	Filters GetCloseDeferredPaymentRequestsFilter
}

func NewGetCloseDeferredPaymentRequests(r *http.Request) (request GetCloseDeferredPaymentRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeCloseDeferredPaymentRequests,
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
