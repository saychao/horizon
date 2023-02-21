package horizon

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
	"github.com/go-errors/errors"
)

type OffersAction struct {
	Action

	accountID         string
	baseAsset         string
	quoteAsset        string
	isBuy             *bool
	offerID           string
	orderBookID       *uint64
	onlyPrimaryMarket bool

	coreRecords []core.Offer
	page        hal.Page
}

// JSON is a method for actions.JSON
func (action *OffersAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.page)
		},
	)
}

func (action *OffersAction) loadParams() {
	action.accountID = action.GetNonEmptyString("account_id")
	action.baseAsset = action.GetString("base_asset")
	action.quoteAsset = action.GetString("quote_asset")
	action.isBuy = action.GetOptionalBool("is_buy")
	action.offerID = action.GetString("offer_id")
	action.orderBookID = action.GetOptionalUint64("order_book_id")
	if (action.baseAsset == "") != (action.quoteAsset == "") {
		action.SetInvalidField("base_asset", errors.New("base and quote assets must be both set or both not set"))
		return
	}

	action.onlyPrimaryMarket = action.GetBool("only_primary")
}

func (action *OffersAction) checkAllowed() {
	action.IsAllowed(action.accountID)
}

func (action *OffersAction) loadRecords() {
	q := action.CoreQ().Offers().ForAccount(action.accountID)
	if action.baseAsset != "" {
		q = q.ForAssets(action.baseAsset, action.quoteAsset)
	}

	if action.isBuy != nil {
		q = q.IsBuy(*action.isBuy)
	}

	if action.offerID != "" {
		q = q.ForOfferID(action.offerID)
	}

	if action.orderBookID != nil {
		q = q.ForOrderBookID(*action.orderBookID)
	}

	if action.onlyPrimaryMarket {
		q = q.OnlyPrimaryMarket()
	}

	err := q.Select(&action.coreRecords)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get offers from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.page.Init()
	for i := range action.coreRecords {
		action.page.Add(resource.PopulateOffer(action.coreRecords[i]))
	}
}
