package horizon

import (
	"github.com/SafeRE-IT/horizon/db2"
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
)

type BalanceIndexAction struct {
	Action
	AccountFilter  string
	ExchangeFilter string
	Asset          string
	PagingParams   db2.PageQuery
	Records        []history.Balance
	Page           hal.Page
}

func (action *BalanceIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *BalanceIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.AccountFilter = action.GetString("account")
	action.ExchangeFilter = action.GetString("exchange")
	action.PagingParams = action.GetPageQuery()
	action.Asset = action.GetString("asset")
	action.Page.Filters = map[string]string{
		"account":  action.AccountFilter,
		"exchange": action.ExchangeFilter,
		"asset":    action.Asset,
	}

}

func (action *BalanceIndexAction) loadRecords() {
	balances := action.HistoryQ().Balances()

	if action.AccountFilter != "" {
		balances.ForAccount(action.AccountFilter)
	}

	if action.Asset != "" {
		balances.ForAsset(action.Asset)
	}

	err := balances.Page(action.PagingParams).Select(&action.Records)

	if err != nil {
		action.Log.WithError(err).Error("failed to get balances")
		action.Err = &problem.ServerError
		return
	}
}

func (action *BalanceIndexAction) loadPage() {
	for _, record := range action.Records {
		var balance resource.BalancePublic
		balance.Populate(record)
		action.Page.Add(balance)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
