package requests

import (
	"net/http"
)

const (
	//IncludeTypeRoleRules - defines if rules of the role should be included into response
	IncludeTypeRoleRules = "rules"
)

var includeTypeAccountRoleAll = map[string]struct{}{
	IncludeTypeRoleRules: {},
}

// GetAccountRole - represents params to be specified by user for Get account role handler
type GetAccountRole struct {
	*base
	ID       uint64
	Includes struct {
		Rules bool `include:"rules"`
	}
}

// NewGetAccountRole returns new instance of GetAsset request
func NewGetAccountRole(r *http.Request) (*GetAccountRole, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAccountRoleAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		// error is already correctly wrapped
		return nil, err
	}

	return &GetAccountRole{
		base: b,
		ID:   id,
	}, nil
}
