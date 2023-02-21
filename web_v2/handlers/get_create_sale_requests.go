package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/db2/core2"

	"gitlab.com/tokend/go/xdr"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetCreateSaleRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateSaleRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateSaleRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	if !isAllowed(r, w, request.GetRequestsBase.Filters.Requestor, request.GetRequestsBase.Filters.Reviewer) {
		return
	}

	err = handler.MakeAll(w, request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get request list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

type getCreateSaleRequestsHandler struct {
	R         requests.GetCreateSaleRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateSaleRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateSaleRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateSale))

	if request.Filters.BaseAsset != nil {
		q = q.FilterBySaleBaseAsset(*request.Filters.BaseAsset)
	}

	if request.Filters.DefaultQuoteAsset != nil {
		q = q.FilterBySaleQuoteAsset(*request.Filters.DefaultQuoteAsset)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateSaleRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateSaleRequestsBaseAsset) {
		base, err := h.AssetsQ.GetByCode(record.Details.CreateSale.BaseAsset)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get base asset")
		}

		if base == nil {
			return regources.ReviewableRequest{}, errors.New("base asset not found")
		}
		resource := resources.NewAsset(*base)
		included.Add(&resource)
	}

	if h.R.ShouldInclude(requests.IncludeTypeCreateSaleRequestsDefaultQuoteAsset) {
		quote, err := h.AssetsQ.GetByCode(record.Details.CreateSale.DefaultQuoteAsset)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get default quote asset")
		}

		if quote == nil {
			return regources.ReviewableRequest{}, errors.New("default quote asset not found")
		}

		resource := resources.NewAsset(*quote)
		included.Add(&resource)
	}

	if h.R.ShouldInclude(requests.IncludeTypeCreateSaleRequestsQuoteAssets) {
		assetCodes := make([]string, 0, len(record.Details.CreateSale.QuoteAssets))
		for _, v := range record.Details.CreateSale.QuoteAssets {
			assetCodes = append(assetCodes, v.Asset)
		}
		quote, err := h.AssetsQ.FilterByCodes(assetCodes).Select()
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get quote assets")
		}

		if quote == nil {
			return regources.ReviewableRequest{}, errors.New("quote assets not found")
		}
		for _, v := range quote {
			resource := resources.NewAsset(v)
			included.Add(&resource)
		}
	}

	return resource, nil
}
