package handlers

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/history2"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetSaleListForAccount - processes request to get the list of sales
func GetSaleListForAccount(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleListForAccountHandler{
		salesBaseHandler: salesBaseHandler{
			SalesQ:           history2.NewSalesQ(historyRepo),
			AssetsQ:          history2.NewAssetQ(historyRepo),
			saleCapConverter: converter,
			Log:              ctx.Log(r),

			ParticipationQ: history2.NewSaleParticipationQ(historyRepo),
			OffersQ:        core2.NewOffersQ(ctx.CoreRepo(r)),
		},
	}

	request, err := requests.NewGetSaleListForAccount(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, &request.Address) {
		return
	}

	result, err := handler.GetSaleListForAccount(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSaleListForAccountHandler struct {
	salesBaseHandler
}

// GetSaleListForAccount returns the list of sales with related resources
func (h *getSaleListForAccountHandler) GetSaleListForAccount(request *requests.GetSaleListForAccount,
) (*regources.SaleListResponse, error) {
	q := applySaleFilters(request.SalesBase, h.SalesQ).Whitelisted(request.Address)
	q = applySaleIncludes(request.SalesBase, q)

	historySales, err := q.CursorPage(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sales list")
	}
	response := &regources.SaleListResponse{
		Data: make([]regources.Sale, 0, len(historySales)),
	}

	err = h.populateResponse(historySales, request.SalesBase, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate response")
	}

	h.populateLinks(response, request)

	return response, nil
}

func (h *getSaleListForAccountHandler) populateLinks(
	response *regources.SaleListResponse, request *requests.GetSaleListForAccount,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}
