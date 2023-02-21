package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	FilterTypeDeferredPaymentListAsset         = "asset"
	FilterTypeDeferredPaymentListSource        = "source"
	FilterTypeDeferredPaymentListSourceBalance = "source_balance"
	FilterTypeDeferredPaymentListDestination   = "destination"
)

var includeTypeDeferredPaymentAll = map[string]struct{}{}

var filterTypeDeferredPaymentListAll = map[string]struct{}{
	FilterTypeDeferredPaymentListAsset:         {},
	FilterTypeDeferredPaymentListSource:        {},
	FilterTypeDeferredPaymentListSourceBalance: {},
	FilterTypeDeferredPaymentListDestination:   {},
}

type GetDeferredPaymentListFilters struct {
	Destination   string `fig:"destination"`
	SourceBalance string `fig:"source_balance"`
	Source        string `fig:"source"`
	Asset         string `fig:"asset"`
}

//GetDeferredPaymentList - represents params to be specified for Get Fees handler
type GetDeferredPaymentList struct {
	*base
	Filters    GetDeferredPaymentListFilters
	PageParams *pgdb.CursorPageParams
}

// NewGetDeferredPaymentList returns the new instance of GetDeferredPaymentList request
func NewGetDeferredPaymentList(r *http.Request) (*GetDeferredPaymentList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeDeferredPaymentListAll,
		supportedIncludes: includeTypeDeferredPaymentAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get page params`")
	}

	request := GetDeferredPaymentList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate filters")
	}

	return &request, nil
}
