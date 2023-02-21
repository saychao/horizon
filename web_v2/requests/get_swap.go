package requests

import "net/http"

const (
	IncludeTypeSwapSourceBalance      = "source_balance"
	IncludeTypeSwapDestinationBalance = "destination_balance"
	IncludeTypeSwapAsset              = "asset"
)

var includeTypeSwapAll = map[string]struct{}{
	IncludeTypeSwapSourceBalance:      {},
	IncludeTypeSwapDestinationBalance: {},
	IncludeTypeSwapAsset:              {},
}

type GetSwap struct {
	*base
	ID       int64
	Includes struct {
		SourceBalance      bool `include:"source_balance"`
		DestinationBalance bool `include:"destination_balance"`
		Asset              bool `include:"asset"`
	}
}

func NewGetSwap(r *http.Request) (*GetSwap, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSwapAll,
		supportedFilters:  map[string]struct{}{},
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}
	return &GetSwap{
		base: b,
		ID:   int64(id),
	}, nil
}
