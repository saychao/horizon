package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeSignerRole - defines if signer roles should be included in the response
	IncludeTypeSignerRole = "role"
	// IncludeTypeSignerRoleRules - defines if rules of signer roles should be included in the response
	IncludeTypeSignerRoleRules = "role.rules"
)

var includeTypeSignerAll = map[string]struct{}{
	IncludeTypeSignerRole:      {},
	IncludeTypeSignerRoleRules: {},
}

//GetAccountSigners - represents params to be specified by user for Get Account Signers handler
type GetAccountSigners struct {
	*base
	Address    string
	PageParams pgdb.OffsetPageParams
	Includes   struct {
		Role      bool `include:"role"`
		RoleRules bool `include:"role.rules"`
	}
}

//NewGetAccountSigners - returns new instance of GetAccountSigners request
func NewGetAccountSigners(r *http.Request) (*GetAccountSigners, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSignerAll,
	})
	if err != nil {
		return nil, err
	}
	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	request := GetAccountSigners{
		base:    b,
		Address: address,
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
