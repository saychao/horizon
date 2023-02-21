package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	FilterTypeSwapListSource             = "source"
	FilterTypeSwapListDestination        = "destination"
	FilterTypeSwapListSourceBalance      = "source_balance"
	FilterTypeSwapListDestinationBalance = "destination_balance"
	FilterTypeSwapListAsset              = "asset"
	FilterTypeSwapListState              = "state"
)

var filterTypeSwapListAll = map[string]struct{}{
	FilterTypeSwapListSource:             {},
	FilterTypeSwapListDestination:        {},
	FilterTypeSwapListSourceBalance:      {},
	FilterTypeSwapListDestinationBalance: {},
	FilterTypeSwapListAsset:              {},
	FilterTypeSwapListState:              {},
}

type GetSwapList struct {
	*base
	Filters struct {
		Source             *string `filter:"source" json:"source"`
		SourceBalance      *string `filter:"source_balance" json:"source_balance"`
		Destination        *string `filter:"destination" json:"destination"`
		DestinationBalance *string `filter:"destination_balance" json:"destination_balance"`
		Asset              *string `filter:"asset" json:"asset"`
		State              *int32  `filter:"state" json:"state"`
	}
	PageParams pgdb.CursorPageParams
}

func NewGetSwapList(r *http.Request) (*GetSwapList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeSwapListAll,
		supportedIncludes: includeTypeSwapAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetSwapList{
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
