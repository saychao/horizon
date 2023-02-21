package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	FilterTypeSaleListParticipant = "participant"
)

var filterTypeSaleListAllWithParticipant = map[string]struct{}{
	FilterTypeSaleListOwner:        {},
	FilterTypeSaleListBaseAsset:    {},
	FilterTypeSaleListAssetType:    {},
	FilterTypeSaleListMaxEndTime:   {},
	FilterTypeSaleListMaxStartTime: {},
	FilterTypeSaleListMinStartTime: {},
	FilterTypeSaleListMinEndTime:   {},
	FilterTypeSaleListState:        {},
	FilterTypeSaleListSaleType:     {},
	FilterTypeSaleListMinHardCap:   {},
	FilterTypeSaleListMinSoftCap:   {},
	FilterTypeSaleListMaxHardCap:   {},
	FilterTypeSaleListMaxSoftCap:   {},
	FilterTypeSaleListParticipant:  {},
	FilterTypeSaleCompanyID:        {},
}

// GetSaleList - represents params to be specified by user for getSaleList handler
type GetSaleList struct {
	SalesBase
	SpecialFilters struct {
		Participant string `json:"participant"`
	}
	PageParams *pgdb.OffsetPageParams
}

// NewGetSaleList returns new instance of GetSaleList request
func NewGetSaleList(r *http.Request) (*GetSaleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleListAll,
		supportedFilters:  filterTypeSaleListAllWithParticipant,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSaleList{
		SalesBase: SalesBase{
			base: b,
		},
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	err = b.populateFilters(&request.SpecialFilters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (g GetSaleList) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"participant": g.SpecialFilters.Participant,
	}
}
