package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeAssetListOwners - defines if the asset owners should be included in the response
	IncludeTypeAssetListOwners = "owner"

	// FilterTypeAssetListOwner - defines if we need to filter the list by owner
	FilterTypeAssetListOwner = "owner"
	// FilterTypeAssetListPolicy - defines if we need to filter the list by policy
	FilterTypeAssetListPolicy = "policy"
	//FilterTypeAssetListState - defines if we need to filter the list by asset state
	FilterTypeAssetListState = "state"
	//FilterTypeAssetListCodes - defines if we need to filter the list by asset codes
	FilterTypeAssetListCodes = "codes"
	// FilterTypeAssetListTypes - defines if we need to filter the list by types
	FilterTypeAssetListTypes = "types"
)

var includeTypeAssetListAll = map[string]struct{}{
	IncludeTypeAssetListOwners: {},
}

var filterTypeAssetListAll = map[string]struct{}{
	FilterTypeAssetListOwner:  {},
	FilterTypeAssetListPolicy: {},
	FilterTypeAssetListState:  {},
	FilterTypeAssetListCodes:  {},
	FilterTypeAssetListTypes:  {},
}

//GetAssetList - represents params to be specified for Get Assets handler
type GetAssetList struct {
	*base

	Policy *uint64  `filter:"policy"`
	Owner  *string  `filter:"owner" default:""`
	State  *uint32  `filter:"state"`
	Types  []uint64 `filter:"types"`
	Codes  []string `filter:"codes"`

	PageParams pgdb.OffsetPageParams

	Includes struct {
		Owner bool `include:"owner"`
	}
}

// NewGetAssetList returns the new instance of GetAssetList request
func NewGetAssetList(r *http.Request) (*GetAssetList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAssetListAll,
		supportedFilters:  filterTypeAssetListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetAssetList{
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
