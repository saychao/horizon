package requests

import "net/http"

// GetSaleForAccount represents params to be specified by user for getSaleForAc handler
type GetSaleForAccount struct {
	GetSale
	Address string
}

// NewGetSale returns new instance of the GetSale request
func NewGetSaleForAccount(r *http.Request) (*GetSaleForAccount, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("sale_id")
	if err != nil {
		return nil, err
	}

	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	return &GetSaleForAccount{
		GetSale: GetSale{base: b,
			ID: id,
		},
		Address: address,
	}, nil
}
