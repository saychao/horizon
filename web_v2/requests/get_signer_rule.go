package requests

import (
	"net/http"
)

// GetSignerRule - represents params to be specified by user for Get signer rule handler
type GetSignerRule struct {
	*base
	ID uint64
}

// NewGetSignerRule returns new instance of GetAsset request
func NewGetSignerRule(r *http.Request) (*GetSignerRule, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		// error is already correctly wrapped
		return nil, err
	}

	return &GetSignerRule{
		base: b,
		ID:   id,
	}, nil
}
