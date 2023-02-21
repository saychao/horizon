package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/history2"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetSaleList - processes request to get the list of sales
func GetSaleList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleListHandler{
		salesBaseHandler: salesBaseHandler{
			SalesQ:           history2.NewSalesQ(historyRepo),
			AssetsQ:          history2.NewAssetQ(historyRepo),
			saleCapConverter: converter,
			Log:              ctx.Log(r),

			ParticipationQ: history2.NewSaleParticipationQ(historyRepo),
			OffersQ:        core2.NewOffersQ(coreRepo),
		},
		OffersQ: core2.NewOffersQ(coreRepo),
	}

	request, err := requests.NewGetSaleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSaleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSaleListHandler struct {
	salesBaseHandler
	OffersQ core2.OffersQ
}

// GetSaleList returns the list of assets with related resources
func (h *getSaleListHandler) GetSaleList(request *requests.GetSaleList) (*regources.SaleListResponse, error) {
	q := applySaleFilters(request.SalesBase, h.SalesQ)
	var err error
	q, err = applyParticipantFilter(request, q, h.OffersQ)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply participant filter")
	}

	q = applySaleIncludes(request.SalesBase, q)

	historySales, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sale list")
	}

	response := &regources.SaleListResponse{
		Data:  make([]regources.Sale, 0, len(historySales)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	err = h.populateResponse(historySales, request.SalesBase, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate response")
	}

	return response, nil
}

func applyParticipantFilter(s *requests.GetSaleList, q history2.SalesQ, offerQ core2.OffersQ,
) (history2.SalesQ, error) {
	if s.ShouldFilter(requests.FilterTypeSaleListParticipant) {
		orderBookIDs, err := offerQ.OrderBookID().FilterByOrderBookID(-1).
			FilterByOwnerID(s.SpecialFilters.Participant).SelectID()
		if err != nil {
			return q, errors.Wrap(err, "failed to select sale ids")
		}

		q = q.FilterByParticipant(s.SpecialFilters.Participant, orderBookIDs)
	}

	return q, nil
}
