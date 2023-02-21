package requests

import "net/http"

type GetBalancesStatistic struct {
	*base
	AssetCode      string
	AccountAddress string
}

func NewGetBalancesStatistic(r *http.Request) (*GetBalancesStatistic, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	accountAddress := b.getString("id")
	assetCode := b.getString("asset_code")

	request := &GetBalancesStatistic{
		base:           b,
		AccountAddress: accountAddress,
		AssetCode:      assetCode,
	}

	return request, nil
}
