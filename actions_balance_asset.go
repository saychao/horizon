package horizon

import (
	"database/sql"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
)

type BalanceAssetAction struct {
	Action
	BalanceID string
	Asset     resource.Asset
}

func (action *BalanceAssetAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Asset)
		},
	)
}

func (action *BalanceAssetAction) loadParams() {
	action.BalanceID = action.GetBalanceIDAsString("balance_id")
}

func (action *BalanceAssetAction) loadData() {
	var result core.Balance
	err := action.CoreQ().BalanceByID(&result, action.BalanceID)
	if err != nil {
		if err == sql.ErrNoRows {
			action.Err = &problem.NotFound
			return
		}

		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("Failed to get balance by id")
		return
	}

	action.Asset.Code = result.Asset
}
