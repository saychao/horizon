package reviewablerequest

import (
	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateSaleCreationRequest(histRequest history.SaleRequest) (
	r *regources.SaleCreationRequest, err error,
) {
	r = &regources.SaleCreationRequest{}
	r.BaseAsset = histRequest.BaseAsset
	r.DefaultQuoteAsset = histRequest.DefaultQuoteAsset
	r.StartTime = regources.Time(histRequest.StartTime)
	r.EndTime = regources.Time(histRequest.EndTime)
	r.SoftCap = histRequest.SoftCap
	r.HardCap = histRequest.HardCap
	r.Details = histRequest.Details
	r.QuoteAssets = histRequest.QuoteAssets
	r.SaleType.Value = int32(histRequest.SaleType)
	r.SaleType.Name = histRequest.SaleType.ShortString()
	r.BaseAssetForHardCap = histRequest.BaseAssetForHardCap
	return
}
