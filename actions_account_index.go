package horizon

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
)

type AccountIndexAction struct {
	Action
	Types []uint64

	Records []core.Account
	Page    hal.Page
}

func (action *AccountIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *AccountIndexAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *AccountIndexAction) loadRecords() {
	// pagination is turned off intentionally, coz we can't have string cursors atm
	err := action.CoreQ().Accounts().
		ForRoles(action.Types).
		Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to load accounts")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountIndexAction) loadPage() {
	for _, record := range action.Records {
		var r resource.Account
		r.Populate(action.Ctx, record)
		action.Page.Add(r)
	}
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
}
