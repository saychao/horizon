package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// FilterTypeAccountListAccount - defines if we need to filter list by account id
	FilterTypeAccountListAccount = "account"
	// FilterTypeAccountListRole - defines if we need to filter list by role
	FilterTypeAccountListRole = "role"
	// IncludeTypeAccountSigners - defines if we need to include by
	IncludeTypeAccountSigners = "signers"
)

var filterTypeAccountListAll = map[string]struct{}{
	FilterTypeAccountListAccount: {},
	FilterTypeAccountListRole:    {},
}

var includeTypeAccountListAll = map[string]struct{}{
	IncludeTypeAccountSigners: {},
}

// GetAccountList - represents params to be specified by user for Get Account list handler
type GetAccountList struct {
	*base
	Filters struct {
		Account []string `filter:"account"`
		Role    []uint64 `filter:"role"`
	}
	Includes struct {
		Signers bool `include:"signers"`
	}
	PageParams pgdb.OffsetPageParams
}

// NewGetAccountList - returns new instance of GetAccountList request
func NewGetAccountList(r *http.Request) (*GetAccountList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeAccountListAll,
		supportedIncludes: includeTypeAccountListAll,
	})
	if err != nil {
		return nil, err
	}

	var request = GetAccountList{
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
