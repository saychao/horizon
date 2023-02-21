package horizon

import (
	"strconv"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
)

type OrderBookAction struct {
	Action
	OwnerID     string
	BaseAsset   string
	QuoteAsset  string
	IsBuy       bool
	OrderBookID *uint64

	CoreRecords []core.Offer
	Page        hal.Page
}

// JSON is a method for actions.JSON
func (action *OrderBookAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkIsSigned,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *OrderBookAction) loadParams() {
	action.BaseAsset = action.GetNonEmptyString("base_asset")
	action.QuoteAsset = action.GetNonEmptyString("quote_asset")
	action.OrderBookID = action.GetOptionalUint64("order_book_id")
	action.OwnerID = action.GetString("owner_id")
	action.IsBuy = action.GetBool("is_buy")
	action.Page.Filters = map[string]string{
		"base_asset":    action.BaseAsset,
		"quote_asset":   action.QuoteAsset,
		"is_buy":        strconv.FormatBool(action.IsBuy),
		"order_book_id": action.GetString("order_book_id"),
	}
}

func (action *OrderBookAction) checkIsSigned() {
	if action.OwnerID != "" {
		action.IsAllowed(action.OwnerID)
		return
	}

	if action.Signer != "" {
		action.isAllowed(action.Signer)
	}
}

func (action *OrderBookAction) loadRecords() {
	q := action.CoreQ().
		Offers().
		ForAssets(action.BaseAsset, action.QuoteAsset).
		IsBuy(action.IsBuy)

	if action.OrderBookID != nil {
		q = q.ForOrderBookID(*action.OrderBookID)
	}

	err := q.OrderByPrice(action.IsBuy).Select(&action.CoreRecords)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get offers from core DB")
		action.Err = &problem.ServerError
		return
	}

	for i := range action.CoreRecords {
		var result resource.OrderBookEntry
		result.Populate(&action.CoreRecords[i].OrderBookEntry, action.BaseAsset, action.QuoteAsset, action.IsBuy)
		if action.IsAdmin || action.OwnerID == action.CoreRecords[i].OwnerID {
			result.OfferID = action.CoreRecords[i].OfferID
			result.OwnerID = action.CoreRecords[i].OwnerID
		}
		action.Page.Add(&result)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.PopulateLinks()
}
