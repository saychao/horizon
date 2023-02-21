package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeOperationsListOperationDetails = "operation.details"

	FilterTypeOperationsListTypes = "types"
	FilterTypeOperationsSource    = "source"
)

type GetOperations struct {
	*base
	PageParams pgdb.CursorPageParams
	Filters    struct {
		Types  []int   `filter:"types"`
		Source *string `filter:"source"`
	}
	Includes struct {
		OperationDetails bool `include:"operation.details"`
	}
}

func NewGetOperations(r *http.Request) (*GetOperations, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeOperationsListOperationDetails: {},
		},
		supportedFilters: map[string]struct{}{
			FilterTypeOperationsListTypes: {},
			FilterTypeOperationsSource:    {},
		},
	})
	if err != nil {
		return nil, err
	}

	request := GetOperations{
		base: b,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = SetDefaultCursorPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
