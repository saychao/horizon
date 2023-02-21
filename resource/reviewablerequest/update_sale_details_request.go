package reviewablerequest

import (
	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateUpdateSaleDetailsRequest(histRequest history.UpdateSaleDetailsRequest) (
	r *regources.UpdateSaleDetailsRequest, err error,
) {
	r = &regources.UpdateSaleDetailsRequest{}
	r.SaleID = histRequest.SaleID
	r.NewDetails = histRequest.NewDetails
	return
}
