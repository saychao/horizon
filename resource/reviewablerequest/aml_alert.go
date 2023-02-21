package reviewablerequest

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	amount2 "gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

func PopulateAmlAlertRequest(histRequest history.AmlAlertRequest) (
	*regources.AMLAlertRequest, error,
) {
	amount := amount2.MustParse(histRequest.Amount)
	return &regources.AMLAlertRequest{
		BalanceID: histRequest.BalanceID,
		Amount:    regources.Amount(amount),
		Reason:    histRequest.Reason,
	}, nil
}
