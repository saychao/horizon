package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"

	"gitlab.com/distributed_lab/kit/pgdb"
)

var includeTypeAccountRoleListAll = map[string]struct{}{
	IncludeTypeRoleRules: {},
}

//GetAccountRoleList - represents params to be specified for Get AccountRoles handler
type GetAccountRoleList struct {
	*base
	PageParams pgdb.OffsetPageParams
	Includes   struct {
		Rules bool `include:"rules"`
	}
}

// NewGetAccountRoleList returns the new instance of GetAccountRoleList request
func NewGetAccountRoleList(r *http.Request) (*GetAccountRoleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAccountRoleListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetAccountRoleList{
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
