package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/kit/pgdb"
	regources "gitlab.com/tokend/regources/generated"
)

const (
	IncludeTypeFeeListAsset = "asset"

	FilterTypeFeeListAsset       = "asset"
	FilterTypeFeeListSubtype     = "subtype"
	FilterTypeFeeListAccount     = "account"
	FilterTypeFeeListAccountRole = "account_role"
	FilterTypeFeeListFeeType     = "fee_type"
	FilterTypeFeeListLowerBound  = "lower_bound"
	FilterTypeFeeListUpperBound  = "upper_bound"
)

var includeTypeFeeListAll = map[string]struct{}{
	IncludeTypeFeeListAsset: {},
}

var filterTypeFeeListAll = map[string]struct{}{
	FilterTypeFeeListAsset:       {},
	FilterTypeFeeListSubtype:     {},
	FilterTypeFeeListAccount:     {},
	FilterTypeFeeListAccountRole: {},
	FilterTypeFeeListLowerBound:  {},
	FilterTypeFeeListUpperBound:  {},
	FilterTypeFeeListFeeType:     {},
}

//GetFeeList - represents params to be specified for Get Fees handler
type GetFeeList struct {
	*base
	Filters    GetFeeListFilters
	PageParams *pgdb.OffsetPageParams
	Includes   struct {
		Asset bool `include:"asset"`
	}
}

type GetFeeListFilters struct {
	Asset       string           `fig:"asset"`
	Subtype     int64            `fig:"subtype"`
	FeeType     int32            `fig:"fee_type"`
	Account     string           `fig:"account"`
	AccountRole uint64           `fig:"account_role"`
	LowerBound  regources.Amount `fig:"lower_bound"`
	UpperBound  regources.Amount `fig:"upper_bound"`
}

// NewGetFeeList returns the new instance of GetFeeList request
func NewGetFeeList(r *http.Request) (*GetFeeList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeFeeListAll,
		supportedFilters:  filterTypeFeeListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get page params`")
	}

	request := GetFeeList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate filters")
	}

	return &request, nil
}
