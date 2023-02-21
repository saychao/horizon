package horizon

import (
	"github.com/SafeRE-IT/horizon/db2"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
)

type HistoryOfferIndexAction struct {
	Action
	OwnerID      string
	BaseAsset    string
	QuoteAsset   string
	PagingParams db2.PageQuery
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *HistoryOfferIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *HistoryOfferIndexAction) checkAllowed() {
	action.IsAllowed(action.OwnerID)
}

func (action *HistoryOfferIndexAction) loadParams() {
	action.OwnerID = action.GetNonEmptyString("owner_id")
	action.BaseAsset = action.GetNonEmptyString("base_asset")
	action.QuoteAsset = action.GetNonEmptyString("quote_asset")

	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"owner_id":    action.OwnerID,
		"base_asset":  action.BaseAsset,
		"quote_asset": action.QuoteAsset,
	}
}

func (action *HistoryOfferIndexAction) loadRecords() {
	historyOffers, err := action.HistoryQ().Offers().
		ForOwnerID(action.OwnerID).
		ForBase(action.BaseAsset).
		ForQuote(action.QuoteAsset).
		Page(action.PagingParams).
		Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to select history offers")
		action.Err = &problem.ServerError
		return
	}

	if len(historyOffers) == 0 {
		action.Err = &problem.NotFound
		return
	}

	for _, offer := range historyOffers {
		action.Page.Add(resource.PopulateHistoryOffer(offer))
	}
}

func (action *HistoryOfferIndexAction) loadPage() {
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
