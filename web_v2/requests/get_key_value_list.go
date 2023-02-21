package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
)

type GetKeyValueList struct {
	*base
	PageParams *pgdb.OffsetPageParams
}

func NewGetKeyValueList(r *http.Request) (*GetKeyValueList, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	return &GetKeyValueList{
		base:       b,
		PageParams: pageParams,
	}, nil
}
