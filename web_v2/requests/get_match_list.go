package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	// FilterTypeMatchListBaseAsset - defines if we need to filter the list by base asset
	FilterTypeMatchListBaseAsset = "base_asset"
	// FilterTypeMatchListQuoteAsset - defines if we need to filter the list by quote asset
	FilterTypeMatchListQuoteAsset = "quote_asset"
)

var filterTypeMatchListAll = map[string]struct{}{
	FilterTypeMatchListBaseAsset:  {},
	FilterTypeMatchListQuoteAsset: {},
}

// GetMatchList represents params to be specified by user for getMatchList handler
type GetMatchList struct {
	*base
	Filters struct {
		BaseAsset  string `fig:"base_asset,required"`
		QuoteAsset string `fig:"quote_asset,required"`
	}

	PageParams *pgdb.CursorPageParams
}

// NewGetMatchList - returns new instance of GetMatchList
func NewGetMatchList(r *http.Request) (*GetMatchList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters: filterTypeMatchListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetMatchList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
