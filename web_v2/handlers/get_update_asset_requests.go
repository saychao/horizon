package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/go/xdr"

	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetUpdateAssetRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUpdateAssetRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getUpdateAssetRequestsHandler{
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

type getUpdateAssetRequestsHandler struct {
	R         requests.GetUpdateAssetRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getUpdateAssetRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetUpdateAssetRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeUpdateAsset))

	if request.Filters.Asset != nil {
		q = q.FilterByAssetUpdateAsset(*request.Filters.Asset)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getUpdateAssetRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeUpdateAssetRequestsAsset) {
		asset, err := h.AssetsQ.GetByCode(record.Details.UpdateAsset.Asset)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get asset")
		}

		if asset == nil {
			return regources.ReviewableRequest{}, errors.New("asset not found")
		}
		resource := resources.NewAsset(*asset)
		included.Add(&resource)
	}

	return resource, nil
}
