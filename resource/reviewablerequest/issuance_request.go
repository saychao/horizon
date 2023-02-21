package reviewablerequest

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

func PopulateIssuanceRequest(histRequest history.IssuanceRequest) (
	r *regources.IssuanceRequest, err error,
) {
	r = &regources.IssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = regources.Amount(amount.MustParse(histRequest.Amount))
	r.Receiver = histRequest.Receiver
	r.ExternalDetails = histRequest.ExternalDetails
	return
}
