package handlers

import (
	"net/http"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// DeprecatedGetOrderBook - processes request to get order book entries
func DeprecatedGetOrderBook(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)

	handler := deprecatedGetOrderBookHandler{
		OrderBooksQ: core2.NewOrderBooksQ(coreRepo),
		SalesQ:      history2.NewSalesQ(historyRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewDeprecatedGetOrderBook(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.DeprecatedGetOrderBook(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get order book entries", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type deprecatedGetOrderBookHandler struct {
	OrderBooksQ core2.OrderBooksQ
	SalesQ      history2.SalesQ
	Log         *logan.Entry
}

const secondaryMarketOrderBookID = 0

// DeprecatedGetOrderBook returns offer with related resources
func (h *deprecatedGetOrderBookHandler) DeprecatedGetOrderBook(request *requests.DeprecatedGetOrderBook) (*regources.OrderBookEntryListResponse, error) {
	if request.ID != secondaryMarketOrderBookID {
		sale, err := h.SalesQ.GetByID(request.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get sale by ID")
		}
		if sale == nil {
			return nil, nil
		}
	}

	q := h.OrderBooksQ.Page(request.PageParams).FilterByOrderBookID(request.ID)

	if request.ShouldInclude(requests.DeprecatedIncludeTypeOrderBookBaseAssets) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.DeprecatedIncludeTypeOrderBookQuoteAssets) {
		q = q.WithQuoteAsset()
	}

	if request.Filters.BaseAsset != nil {
		q = q.FilterByBaseAssetCode(*request.Filters.BaseAsset)
	}

	if request.Filters.QuoteAsset != nil {
		q = q.FilterByQuoteAssetCode(*request.Filters.QuoteAsset)
	}

	if request.Filters.IsBuy != nil {
		q = q.FilterByIsBuy(*request.Filters.IsBuy)
	}

	coreOrderBookEntries, err := q.Select()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get offer list")
	}

	response := &regources.OrderBookEntryListResponse{
		Data:  make([]regources.OrderBookEntry, 0, len(coreOrderBookEntries)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, coreOrderBookEntry := range coreOrderBookEntries {
		setPoints(&coreOrderBookEntry)

		response.Data = append(response.Data, resources.NewOrderBookEntry(coreOrderBookEntry))

		if request.ShouldInclude(requests.DeprecatedIncludeTypeOrderBookBaseAssets) {
			baseAsset := resources.NewAsset(*coreOrderBookEntry.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.DeprecatedIncludeTypeOrderBookQuoteAssets) {
			quoteAsset := resources.NewAsset(*coreOrderBookEntry.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	return response, nil
}
