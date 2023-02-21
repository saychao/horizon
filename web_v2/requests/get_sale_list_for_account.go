package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
)

// GetSaleList - represents params to be specified by user for getSaleList handler
type GetSaleListForAccount struct {
	SalesBase
	Address    string
	PageParams *pgdb.CursorPageParams
}

// NewGetSaleListForAccount returns new instance of GetSaleListForAccount request
func NewGetSaleListForAccount(r *http.Request) (*GetSaleListForAccount, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleListAll,
		supportedFilters:  filterTypeSaleListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	request := GetSaleListForAccount{
		Address: address,
		SalesBase: SalesBase{
			base: b,
		},
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
