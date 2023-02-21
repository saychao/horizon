package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeKYCRecoveryRequestsAccount = "request_details.target_account"
)

var includeTypeKYCRecoveryRequests = map[string]struct{}{
	IncludeTypeKYCRecoveryRequestsAccount: {},
}

const (
	FilterTypeKYCRecoveryRequestsAccount = "request_details.target_account"
)

var filterTypeKYCRecoveryRequests = map[string]struct{}{
	FilterTypeKYCRecoveryRequestsAccount: {},
}

type GetKYCRecoveryRequestsFilter struct {
	Account *string `filter:"request_details.target_account"`
}

type GetKYCRecoveryRequests struct {
	GetRequestsBase
	Filters  GetKYCRecoveryRequestsFilter
	Includes struct {
		TargetAccount bool `include:"request_details.target_account"`
	}
}

func NewGetKYCRecoveryRequests(r *http.Request) (request GetKYCRecoveryRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		filterTypeKYCRecoveryRequests,
		includeTypeKYCRecoveryRequests,
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
