package horizon

import (
	"github.com/SafeRE-IT/horizon/db2"
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
)

type TradesAction struct {
	Action
	BaseAsset    string
	QuoteAsset   string
	PagingParams db2.PageQuery
	OrderBookID  *uint64

	Trades []history.Trades
	Page   hal.Page
}

// JSON is a method for actions.JSON
func (action *TradesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *TradesAction) loadParams() {
	action.BaseAsset = action.GetNonEmptyString("base_asset")
	action.QuoteAsset = action.GetNonEmptyString("quote_asset")
	action.OrderBookID = action.GetOptionalUint64("order_book_id")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"base_asset":    action.BaseAsset,
		"quote_asset":   action.QuoteAsset,
		"order_book_id": action.GetString("order_book_id"),
	}
}

func (action *TradesAction) loadRecords() {
	q := action.HistoryQ().Trades().ForPair(action.BaseAsset, action.QuoteAsset).Page(action.PagingParams)
	if action.OrderBookID != nil {
		q = q.ForOrderBook(*action.OrderBookID)
	}

	err := q.Select(&action.Trades)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get trades")
		action.Err = &problem.ServerError
		return
	}

	for i := range action.Trades {
		var result resource.Trades
		result.Populate(&action.Trades[i])
		action.Page.Add(&result)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()

}
