package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// FilterTypeSaleWhitelistAddress - defines if we need to filter entries by account address they are applied to
	FilterTypeSaleWhitelistAddress = "address"
)

// GetSaleWhitelist - represents params to be specified by user for getSaleWhitelist handler
type GetSaleWhitelist struct {
	*base
	SaleID  uint64
	Filters struct {
		Address *string `filter:"address"`
	}
	PageParams pgdb.CursorPageParams
}

// NewGetSaleWhitelist returns new instance of GetSaleWhitelist
func NewGetSaleWhitelist(r *http.Request) (*GetSaleWhitelist, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters: map[string]struct{}{
			FilterTypeSaleWhitelistAddress: {},
		},
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	request := GetSaleWhitelist{
		base:   b,
		SaleID: id,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = SetDefaultCursorPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
