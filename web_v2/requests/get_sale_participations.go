package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeSaleParticipationsBaseAsset - defines if the base asset should be included in the response
	IncludeTypeSaleParticipationsBaseAsset = "base_asset"
	// IncludeTypeSaleParticipationsQuoteAsset - defines if the quote asset should be included in the response
	IncludeTypeSaleParticipationsQuoteAsset = "quote_asset"

	// FilterTypeSaleParticipationsParticipant - defines if we need to filter response by participant
	FilterTypeSaleParticipationsParticipant = "participant"
	// FilterTypeSaleParticipationsQuoteAsset - defines if we need to filter response by quote_asset
	FilterTypeSaleParticipationsQuoteAsset = "quote_asset"
)

var includeTypeSaleParticipationAll = map[string]struct{}{
	IncludeTypeSaleParticipationsQuoteAsset: {},
	IncludeTypeSaleParticipationsBaseAsset:  {},
}

var filterTypeSaleParticipationAll = map[string]struct{}{
	FilterTypeSaleParticipationsParticipant: {},
	FilterTypeSaleParticipationsQuoteAsset:  {},
}

// GetSaleParticipations - represents params to be specified by user for getSaleParticipation handler
type GetSaleParticipations struct {
	*base
	SaleID  uint64
	Filters struct {
		QuoteAsset  *string `filter:"quote_asset" json:"quote_asset"`
		Participant *string `filter:"participant" json:"participant"`
	}
	PageParams pgdb.CursorPageParams
	Includes   struct {
		BaseAsset   bool `include:"base_asset"`
		QuoteAssets bool `include:"quote_assets"`
	}
}

// NewGetSaleParticipations returns new instance of GetSaleParticipations
func NewGetSaleParticipations(r *http.Request) (*GetSaleParticipations, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeSaleParticipationAll,
		supportedIncludes: includeTypeSaleParticipationAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	request := GetSaleParticipations{
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
