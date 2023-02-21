package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	// FilterTypeConvertedBalancesAssetOwner - defines if we need to filter the list by asset owner
	FilterTypeConvertedBalancesAssetOwner = "asset_owner"

	// IncludeTypeConvertedBalancesAsset - defines if conversion asset should be included in the response
	IncludeTypeConvertedBalancesAsset = "asset"
	// IncludeTypeConvertedBalancesStates - defines if converted balances states should be included in the response
	IncludeTypeConvertedBalancesStates = "states"
	// IncludeTypeConvertedBalancesBalance - defines if balances states should be included in the response
	IncludeTypeConvertedBalancesBalance = "balance"
	// IncludeTypeConvertedBalancesBalanceState - defines if balances states should be included in the response
	IncludeTypeConvertedBalancesBalanceState = "balance.state"
	// IncludeTypeConvertedBalancesBalanceAsset - defines if balances assets should be included in the response
	IncludeTypeConvertedBalancesBalanceAsset = "balance.asset"
)

var includeTypeConvertedBalancesAll = map[string]struct{}{
	IncludeTypeConvertedBalancesAsset:        {},
	IncludeTypeConvertedBalancesStates:       {},
	IncludeTypeConvertedBalancesBalance:      {},
	IncludeTypeConvertedBalancesBalanceState: {},
	IncludeTypeConvertedBalancesBalanceAsset: {},
}

var filterTypeConvertedBalancesAll = map[string]struct{}{
	FilterTypeBalanceListAssetOwner: {},
}

// GetConvertedBalances - represents params to be specified by user for GetConvertedBalances handler
type GetConvertedBalances struct {
	*base
	Filters struct {
		AssetOwner *string `filter:"asset_owner" json:"asset_owner"`
	}
	Includes struct {
		Asset        bool `include:"asset"`
		States       bool `include:"states"`
		Balance      bool `include:"balance"`
		BalanceState bool `include:"balance.state"`
		BalanceAsset bool `include:"balance.asset"`
	}
	AssetCode      string
	AccountAddress string
}

// NewGetConvertedBalances returns new instance of GetConvertedBalances request
func NewGetConvertedBalances(r *http.Request) (*GetConvertedBalances, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeConvertedBalancesAll,
		supportedFilters:  filterTypeConvertedBalancesAll,
	})
	if err != nil {
		return nil, err
	}

	accountAddress := b.getString("id")
	assetCode := b.getString("asset_code")

	request := GetConvertedBalances{
		base:           b,
		AccountAddress: accountAddress,
		AssetCode:      assetCode,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
