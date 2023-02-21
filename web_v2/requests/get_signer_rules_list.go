package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

//GetSignerRuleList - represents params to be specified for Get SignerRules handler
type GetSignerRuleList struct {
	*base
	PageParams pgdb.OffsetPageParams
}

// NewGetSignerRuleList returns the new instance of GetSignerRuleList request
func NewGetSignerRuleList(r *http.Request) (*GetSignerRuleList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	request := GetSignerRuleList{
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
