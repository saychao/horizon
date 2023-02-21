package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeRedemptionRequestsSourceBalance      = "request_details.source_balance"
	FilterTypeRedemptionRequestsDestinationAccount = "request_details.dest_account"
)

var filterTypeRedemptionRequests = map[string]struct{}{
	FilterTypeRedemptionRequestsSourceBalance:      {},
	FilterTypeRedemptionRequestsDestinationAccount: {},
}

const (
	IncludeTypeRedemptionRequestsSourceBalance      = "request_details.source_balance"
	IncludeTypeRedemptionRequestsDestinationAccount = "request_details.dest_account"
)

var includeTypeRedemptionRequests = map[string]struct{}{
	IncludeTypeRedemptionRequestsSourceBalance:      {},
	IncludeTypeRedemptionRequestsDestinationAccount: {},
}

type GetRedemptionRequestsFilter struct {
	SourceBalance      *string `filter:"request_details.source_balance"`
	DestinationAccount *string `filter:"request_details.dest_account"`
}

type GetRedemptionRequests struct {
	GetRequestsBase
	Filters  GetRedemptionRequestsFilter
	Includes struct {
		RequestDetailsSourceBalance bool `include:"request_details.source_balance"`
		RequestDetailsDestAccount   bool `include:"request_details.dest_account"`
	}
}

func NewGetRedemptionRequests(r *http.Request) (request GetRedemptionRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeRedemptionRequests,
		includeTypeRedemptionRequests,
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
