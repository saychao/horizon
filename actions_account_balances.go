package horizon

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
)

type AccountBalancesAction struct {
	Action

	AccountID string

	Records  []history.Balance
	Resource []resource.BalancePublic
}

func (action *AccountBalancesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountBalancesAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *AccountBalancesAction) loadRecords() {
	if err := action.HistoryQ().Balances().ForAccount(action.AccountID).Select(&action.Records); err != nil {
		action.Log.WithError(err).Error("failed to get balances")
		action.Err = &problem.ServerError
		return
	}

	if len(action.Records) == 0 {
		action.Err = &problem.NotFound
		return
	}
}

func (action *AccountBalancesAction) loadResource() {
	for _, record := range action.Records {
		var r resource.BalancePublic
		r.Populate(record)
		action.Resource = append(action.Resource, r)
	}
}
