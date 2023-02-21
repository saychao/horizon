package requests

import (
	"net/http"
)

// GetAccountRule - represents params to be specified by user for Get account rule handler
type GetAccountRule struct {
	*base
	ID uint64
}

// NewGetAccountRule returns new instance of GetAsset request
func NewGetAccountRule(r *http.Request) (*GetAccountRule, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		// error is already correctly wrapped
		return nil, err
	}

	return &GetAccountRule{
		base: b,
		ID:   id,
	}, nil
}
