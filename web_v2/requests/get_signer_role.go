package requests

import (
	"net/http"
)

var includeTypeSignerRoleAll = map[string]struct{}{
	IncludeTypeRoleRules: {},
}

// GetSignerRole - represents params to be specified by user for Get signer role handler
type GetSignerRole struct {
	*base
	ID uint64
}

// NewGetSignerRole returns new instance of GetAsset request
func NewGetSignerRole(r *http.Request) (*GetSignerRole, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSignerRoleAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		// error is already correctly wrapped
		return nil, err
	}

	return &GetSignerRole{
		base: b,
		ID:   id,
	}, nil
}
