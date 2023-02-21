package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"

	"net/http"
)

const (
	IncludeTypeAskListBaseBalances = "base_balance"
	IncludeTypeAskListOwners       = "owner"
	IncludeTypeAskListBaseAssets   = "base_asset"
	IncludeTypeAskListQuoteAssets  = "quote_assets"

	FilterTypeAskListOwner       = "owner"
	FilterTypeAskListBaseAsset   = "base_asset"
	FilterTypeAskListQuoteAssets = "quote_assets"
)

var includeTypeAskListAll = map[string]struct{}{
	IncludeTypeAskListBaseBalances: {},
	IncludeTypeAskListOwners:       {},
	IncludeTypeAskListBaseAssets:   {},
	IncludeTypeAskListQuoteAssets:  {},
}

var filterTypeAskListAll = map[string]struct{}{
	FilterTypeAskListOwner:       {},
	FilterTypeAskListBaseAsset:   {},
	FilterTypeAskListQuoteAssets: {},
}

// GetAtomicSwapAskList - represents params to be specified by user for Get AtomicSwapAskList handler
type GetAtomicSwapAskList struct {
	*base
	Filters struct {
		Owner       *string  `filter:"owner"`
		BaseAsset   []string `filter:"base_asset"`
		QuoteAssets []string `filter:"quote_assets"`
	}
	PageParams pgdb.OffsetPageParams
	Includes   struct {
		BaseBalance bool `include:"base_balance"`
		Owner       bool `include:"owner"`
		BaseAsset   bool `include:"base_asset"`
		QuoteAssets bool `include:"quote_assets"`
	}
}

// NewGetAtomicSwapAskList returns new instance of GetAtomicSwapAskList request
func NewGetAtomicSwapAskList(r *http.Request) (*GetAtomicSwapAskList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAskListAll,
		supportedFilters:  filterTypeAskListAll,
	})
	if err != nil {
		return nil, err
	}

	// bid relations has not asset relation, we use balance relation
	if _, ok := b.include[IncludeTypeAskBaseAsset]; ok {
		b.include[IncludeTypeAskBaseBalance] = struct{}{}
	}

	request := GetAtomicSwapAskList{
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
