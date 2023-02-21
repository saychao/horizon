package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

var includeTypeSignerRoleListAll = map[string]struct{}{
	IncludeTypeRoleRules: {},
}

//GetSignerRoleList - represents params to be specified for Get SignerRoles handler
type GetSignerRoleList struct {
	*base
	PageParams pgdb.OffsetPageParams
	Includes   struct {
		Rules bool `include:"rules"`
	}
}

// NewGetSignerRoleList returns the new instance of GetSignerRoleList request
func NewGetSignerRoleList(r *http.Request) (*GetSignerRoleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSignerRoleListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetSignerRoleList{
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
