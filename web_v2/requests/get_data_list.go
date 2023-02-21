package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	IncludeTypeDataOwner = "owner"

	FilterTypeDataListOwner = "owner"
	FilterTypeDataListType  = "type"
)

var includeTypeDataAll = map[string]struct{}{
	IncludeTypeDataOwner: {},
}

var filterTypeDataListAll = map[string]struct{}{
	FilterTypeDataListOwner: {},
	FilterTypeDataListType:  {},
}

type GetDataListFilters struct {
	Owner string `fig:"owner"`
	Type  int64  `fig:"type"`
}

//GetDataList - represents params to be specified for Get Fees handler
type GetDataList struct {
	*base
	Filters    GetDataListFilters
	PageParams *pgdb.CursorPageParams
}

// NewGetDataList returns the new instance of GetDataList request
func NewGetDataList(r *http.Request) (*GetDataList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeDataListAll,
		supportedIncludes: includeTypeDataAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get page params`")
	}

	request := GetDataList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate filters")
	}

	return &request, nil
}
