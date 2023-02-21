package reviewablerequest

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateChangeRoleRequest(request *history.ReviewableRequest, histRequest history.ChangeRoleRequest) (
	r *regources.ChangeRoleRequest, err error,
) {
	r = &regources.ChangeRoleRequest{}
	r.DestinationAccount = histRequest.DestinationAccount
	r.AccountRoleToSet = histRequest.AccountRoleToSet
	r.KYCData = histRequest.KYCData
	r.AllTasks = request.AllTasks
	r.PendingTasks = request.PendingTasks
	r.SequenceNumber = histRequest.SequenceNumber
	return
}
