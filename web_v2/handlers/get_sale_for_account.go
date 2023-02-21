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
	regources "gitlab.com/tokend/regources/generated"
)

// GetSaleForAccount - processes request to get sale and it's details by sale ID
func GetSaleForAccount(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleForAccountHandler{
		getSaleBase{
			SalesQ:           history2.NewSalesQ(historyRepo),
			AssetsQ:          history2.NewAssetQ(historyRepo),
			saleCapConverter: converter,
			Log:              ctx.Log(r),
			ParticipationQ:   history2.NewSaleParticipationQ(historyRepo),
			OffersQ:          core2.NewOffersQ(ctx.CoreRepo(r)),
		},
	}

	request, err := requests.NewGetSaleForAccount(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, &request.Address) {
		return
	}

	result, err := handler.GetSale(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale", logan.F{
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

type getSaleForAccountHandler struct {
	getSaleBase
}

// GetSale returns sale with related resources
func (h *getSaleForAccountHandler) GetSale(request *requests.GetSaleForAccount) (*regources.SaleResponse, error) {
	q := h.SalesQ.Whitelisted(request.Address).FilterByID(request.ID)

	if request.ShouldInclude(requests.IncludeTypeSaleListBaseAssets) {
		q = q.WithAsset()
	}

	return h.getAndPopulateResponse(q, &request.GetSale)
}
