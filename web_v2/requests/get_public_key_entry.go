package requests

import (
	"net/http"
)

// GetPublicKeyEntry - represents params to be specified by user for GetPublicKeyEntry handler
type GetPublicKeyEntry struct {
	*base
	ID string
}

// NewGetPublicKeyEntry - returns new instance of GetPublicKeyEntry request
func NewGetPublicKeyEntry(r *http.Request) (*GetPublicKeyEntry, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	id := b.getString("id")

	return &GetPublicKeyEntry{
		base: b,
		ID:   id,
	}, nil
}
