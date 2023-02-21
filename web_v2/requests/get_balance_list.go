package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeBalanceListState - defines if the state of the balance should be included in the response
	IncludeTypeBalanceListState = "state"

	// IncludeTypeBalanceListOwner - defines if the owner of the balance should be included in the response
	IncludeTypeBalanceListOwner = "owner"

	// FilterTypeBalanceListAsset - defines if we need to filter the list by asset
	FilterTypeBalanceListAsset = "asset"

	// FilterTypeBalanceListOwner - defines if we need to filter the list by owner
	FilterTypeBalanceListOwner = "owner"

	// FilterTypeBalanceListAssetOwner - defines if we need to filter the list by asset owner
	FilterTypeBalanceListAssetOwner = "asset_owner"

	//FilterTypeBalanceListAmountLwOrEq - defines if we need to filter the list by amount lower or equal then provided
	FilterTypeBalanceListAmountLwOrEq = "amount_lw_or_eq"

	//FilterTypeBalanceListAmountGt - defines if we need to filter the list by amount greater then provided
	FilterTypeBalanceListAmountGt = "amount_gt"
)

var includeTypeBalanceListAll = map[string]struct{}{
	IncludeTypeBalanceListState: {},
	IncludeTypeBalanceListOwner: {},
}

var filterTypeBalanceListAll = map[string]struct{}{
	FilterTypeBalanceListAsset:        {},
	FilterTypeBalanceListAssetOwner:   {},
	FilterTypeBalanceListOwner:        {},
	FilterTypeBalanceListAmountLwOrEq: {},
	FilterTypeBalanceListAmountGt:     {},
}

// GetBalanceList - represents params to be specified by user for getBalanceList handler
type GetBalanceList struct {
	*base
	Filters struct {
		Asset        *string `filter:"asset"`
		AssetOwner   *string `filter:"asset_owner"`
		Owner        *string `filter:"owner"`
		AmountLwOrEq *uint64 `filter:"amount_lw_or_eq"`
		AmountGt     *uint64 `filter:"amount_gt"`
	}
	Includes struct {
		State bool `include:"state"`
		Owner bool `include:"owner"`
	}
	PageParams pgdb.OffsetPageParams
}

// NewGetBalanceList - returns new instance of GetBalanceList
func NewGetBalanceList(r *http.Request) (*GetBalanceList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeBalanceListAll,
		supportedFilters:  filterTypeBalanceListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetBalanceList{
		base: b,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = b.SetDefaultOffsetPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
