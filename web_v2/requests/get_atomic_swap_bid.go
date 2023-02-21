package requests

import (
	"net/http"
)

const (
	IncludeTypeAskBaseBalance = "base_balance"
	IncludeTypeAskOwner       = "owner"
	IncludeTypeAskBaseAsset   = "base_asset"
	IncludeTypeAskQuoteAssets = "quote_assets"
)

var includeTypeAskAll = map[string]struct{}{
	IncludeTypeAskBaseBalance: {},
	IncludeTypeAskOwner:       {},
	IncludeTypeAskBaseAsset:   {},
	IncludeTypeAskQuoteAssets: {},
}

// GetAtomicSwapAsk - represents params to be specified by user for Get AtomicSwapBid handler
type GetAtomicSwapAsk struct {
	*base
	ID       uint64
	Includes struct {
		BaseBalance bool `include:"base_balance"`
		Owner       bool `include:"owner"`
		BaseAsset   bool `include:"base_asset"`
		QuoteAsset  bool `include:"quote_asset"`
	}
}

// NewGetAtomicSwapAsk returns new instance of GetAtomicSwapAsk request
func NewGetAtomicSwapAsk(r *http.Request) (*GetAtomicSwapAsk, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAskAll,
	})
	if err != nil {
		return nil, err
	}

	// bid relations has not asset relation, we use balance relation
	if _, ok := b.include[IncludeTypeAskBaseAsset]; ok {
		b.include[IncludeTypeAskBaseBalance] = struct{}{}
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetAtomicSwapAsk{
		base: b,
		ID:   id,
	}, nil
}
