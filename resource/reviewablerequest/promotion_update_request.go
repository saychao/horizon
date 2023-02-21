package reviewablerequest

import (
	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulatePromotionUpdateRequest(histRequest history.PromotionUpdateRequest) (
	r *regources.PromotionUpdateRequest, err error,
) {
	r = &regources.PromotionUpdateRequest{}
	r.SaleID = histRequest.SaleID
	newPromotionData, err := PopulateSaleCreationRequest(histRequest.NewPromotionData)
	if newPromotionData != nil {
		r.NewPromotionData = *newPromotionData
	}
	return
}
