package reviewablerequest

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

func PopulatePreIssuanceRequest(histRequest history.PreIssuanceRequest) (
	r *regources.PreIssuanceRequest, err error,
) {
	r = &regources.PreIssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = regources.Amount(amount.MustParse(histRequest.Amount))
	r.Signature = histRequest.Signature
	r.Reference = histRequest.Reference
	return
}
