package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

//GetAccountRuleList - represents params to be specified for Get AccountRules handler
type GetAccountRuleList struct {
	*base
	PageParams pgdb.OffsetPageParams
}

// NewGetAccountRuleList returns the new instance of GetAccountRuleList request
func NewGetAccountRuleList(r *http.Request) (*GetAccountRuleList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	request := GetAccountRuleList{
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
