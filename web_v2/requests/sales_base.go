package requests

import (
	"time"

	regources "gitlab.com/tokend/regources/generated"
)

const (
	IncludeTypeSaleListBaseAssets        = "base_asset"
	IncludeTypeSaleListQuoteAssets       = "quote_assets"
	IncludeTypeSaleListDefaultQuoteAsset = "default_quote_asset"

	FilterTypeSaleListOwner        = "owner"
	FilterTypeSaleListBaseAsset    = "base_asset"
	FilterTypeSaleListAssetType    = "asset_type"
	FilterTypeSaleListMaxEndTime   = "max_end_time"
	FilterTypeSaleListMaxStartTime = "max_start_time"
	FilterTypeSaleListMinStartTime = "min_start_time"
	FilterTypeSaleListMinEndTime   = "min_end_time"
	FilterTypeSaleListState        = "state"
	FilterTypeSaleListSaleType     = "sale_type"
	FilterTypeSaleListMinHardCap   = "min_hard_cap"
	FilterTypeSaleListMinSoftCap   = "min_soft_cap"
	FilterTypeSaleListMaxHardCap   = "max_hard_cap"
	FilterTypeSaleListMaxSoftCap   = "max_soft_cap"
	FilterTypeSaleListIDs          = "ids"
	FilterTypeSaleCompanyID        = "company_id"
)

var includeTypeSaleListAll = map[string]struct{}{
	IncludeTypeSaleListBaseAssets:        {},
	IncludeTypeSaleListQuoteAssets:       {},
	IncludeTypeSaleListDefaultQuoteAsset: {},
}

var filterTypeSaleListAll = map[string]struct{}{
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
	FilterTypeSaleListIDs:          {},
	FilterTypeSaleCompanyID:        {},
}

type SalesBase struct {
	*base
	Filters struct {
		Owner        string           `json:"owner"`
		BaseAsset    string           `json:"base_asset"`
		AssetType    uint64           `json:"asset_type"`
		MaxEndTime   *time.Time       `json:"max_end_time"`
		MaxStartTime *time.Time       `json:"max_start_time"`
		MinStartTime *time.Time       `json:"min_start_time"`
		MinEndTime   *time.Time       `json:"min_end_time"`
		State        uint64           `json:"state"`
		SaleType     uint64           `json:"sale_type"`
		MinHardCap   regources.Amount `json:"min_hard_cap"`
		MinSoftCap   regources.Amount `json:"min_soft_cap"`
		MaxHardCap   regources.Amount `json:"max_hard_cap"`
		MaxSoftCap   regources.Amount `json:"max_soft_cap"`
		IDs          []uint64         `json:"ids" fig:"ids"`
		CompanyIDs   []uint64         `json:"company_id" fig:"company_id"`
	}
	Includes struct {
		BaseAsset         bool `include:"base_asset"`
		QuoteAssets       bool `include:"quote_assets"`
		DefaultQuoteAsset bool `include:"default_quote_asset"`
	}
}
