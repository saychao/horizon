package requests

import (
	"net/http"
)

const (
	// IncludeTypeBalanceListState - defines if the state of the balance should be included in the response
	IncludeTypeBalanceState = "state"
	IncludeTypeBalanceAsset = "asset"
)

var includeTypeBalanceAll = map[string]struct{}{
	IncludeTypeBalanceState: {},
	IncludeTypeBalanceAsset: {},
}

type GetBalance struct {
	*base
	BalanceID string
	Includes  struct {
		State bool `include:"state"`
		Asset bool `include:"asset"`
	}
}

func NewGetBalance(r *http.Request) (*GetBalance, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeBalanceAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := newBalanceID(b, "id")
	if err != nil {
		return nil, err
	}

	return &GetBalance{
		base:      b,
		BalanceID: id,
	}, nil
}
