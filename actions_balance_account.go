package horizon

import (
	"database/sql"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
)

type BalanceAccountAction struct {
	Action

	BalanceID string

	Record history.Account

	Resource resource.HistoryAccount
}

func (action *BalanceAccountAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *BalanceAccountAction) loadParams() {
	action.BalanceID = action.GetNonEmptyString("balance_id")
}

func (action *BalanceAccountAction) loadRecord() {
	var balance core.Balance
	err := action.CoreQ().BalanceByID(&balance, action.BalanceID)
	if err == sql.ErrNoRows {
		action.Err = &problem.NotFound
		return
	}
	if err != nil {
		action.Log.WithError(err).Error("failed to get balance")
		action.Err = &problem.ServerError
		return
	}
	action.Record = history.Account{
		Address: balance.AccountID,
	}
}

func (action *BalanceAccountAction) loadResource() {
	action.Resource.Populate(action.Record)
}
