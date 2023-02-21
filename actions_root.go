package horizon

import (
	"github.com/saychao/horizon/ledger"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/tokend/go/xdr"
)

// RootAction provides a summary of the horizon instance and links to various
// useful endpoints
type RootAction struct {
	Action
}

// JSON renders the json response for RootAction
func (action *RootAction) JSON() {
	action.App.UpdateCoreInfo()

	if action.App.CoreInfo == nil {
		action.Err = &problem.ServerOverCapacity
		return
	}

	var res resource.Root
	res.PopulateLedgerState(action.Ctx, ledger.CurrentState())

	res.NetworkPassphrase = action.App.CoreInfo.NetworkPassphrase
	res.AdminAccountID = action.App.CoreInfo.AdminAccountID
	res.MasterAccountID = action.App.CoreInfo.AdminAccountID
	res.MasterExchangeName = action.App.CoreInfo.MasterExchangeName
	res.EnvironmentName = action.App.CoreInfo.MasterExchangeName
	res.TxExpirationPeriod = action.App.CoreInfo.TxExpirationPeriod
	res.XDRRevision = xdr.Revision
	res.HorizonRevision = action.App.horizonVersion
	hal.Render(action.W, res)
}
