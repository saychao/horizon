package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeLimitsListAsset = "asset"

	FilterTypeLimitsListAsset       = "asset"
	FilterTypeLimitsListAccount     = "account"
	FilterTypeLimitsListAccountRole = "account_role"
	FilterTypeLimitsListStatsOpType = "stats_op_type"
)

var includeTypeLimitsListAll = map[string]struct{}{
	IncludeTypeLimitsListAsset: {},
}

var filterTypeLimitsListAll = map[string]struct{}{
	FilterTypeLimitsListAsset:       {},
	FilterTypeLimitsListAccount:     {},
	FilterTypeLimitsListAccountRole: {},
	FilterTypeLimitsListStatsOpType: {},
}

//GetLimitsList - represents params to be specified for Get Fees handler
type GetLimitsList struct {
	*base
	Filters    GetLimitsListFilters
	PageParams pgdb.OffsetPageParams
	Includes   struct {
		Asset bool `include:"asset"`
	}
}
type GetLimitsListFilters struct {
	Asset       *string `filter:"asset"`
	StatsOpType *int32  `filter:"stats_op_type"`
	Account     *string `filter:"account"`
	AccountRole *uint64 `filter:"account_role"`
}

// NewGetLimitsList returns the new instance of GetLimitsList request
func NewGetLimitsList(r *http.Request) (*GetLimitsList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeLimitsListAll,
		supportedFilters:  filterTypeLimitsListAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetLimitsList{
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
