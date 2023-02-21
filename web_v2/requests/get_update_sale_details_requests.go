package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeUpdateSaleDetailsRequestsSale = "request_details.sale"
)

var includeTypeUpdateSaleDetailsRequests = map[string]struct{}{
	IncludeTypeUpdateSaleDetailsRequestsSale: {},
}

type GetUpdateSaleDetailsRequests struct {
	GetRequestsBase
	Includes struct {
		RequestDetailsSale bool `include:"request_details.sale"`
	}
}

func NewGetUpdateSaleDetailsRequests(r *http.Request) (request GetUpdateSaleDetailsRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		map[string]struct{}{},
		includeTypeUpdateSaleDetailsRequests,
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
