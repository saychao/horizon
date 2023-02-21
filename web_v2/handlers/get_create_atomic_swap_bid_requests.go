package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/go/xdr"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetCreateAtomicSwapBidRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateAtomicSwapBidRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateAtomicSwapBidRequestsHandler{
		R:            request,
		RequestsQ:    history2.NewReviewableRequestsQ(historyRepo),
		AssetsQ:      core2.NewAssetsQ(coreRepo),
		QuoteAssetsQ: core2.NewAtomicSwapQuoteAssetQ(coreRepo),
		AsksQ:        core2.NewAtomicSwapAskQ(coreRepo),
		Log:          ctx.Log(r),
	}

	if !isAllowed(r, w, request.GetRequestsBase.Filters.Requestor,
		request.GetRequestsBase.Filters.Reviewer, request.Filters.AskOwner,
	) {
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

type getCreateAtomicSwapBidRequestsHandler struct {
	R            requests.GetCreateAtomicSwapBidRequests
	Base         getRequestListBaseHandler
	RequestsQ    history2.ReviewableRequestsQ
	AssetsQ      core2.AssetsQ
	QuoteAssetsQ core2.AtomicSwapQuoteAssetQ
	AsksQ        core2.AtomicSwapAskQ
	Log          *logan.Entry
}

func (h *getCreateAtomicSwapBidRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateAtomicSwapBidRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateAtomicSwapBid))
	var ids []uint64

	if request.Filters.QuoteAsset != nil {
		q = q.FilterByAtomicSwapQuoteAsset(*request.Filters.QuoteAsset)
	}
	if request.Filters.AskOwner != nil {
		var err error
		ids, err = h.AsksQ.IDSelector().FilterByOwner(*request.Filters.AskOwner).SelectIDs()
		if err != nil {
			return errors.Wrap(err, "failed to load ask ids by owner", logan.F{
				"owner": *request.Filters.AskOwner,
			})
		}

		q = q.FilterByAtomicSwapAskIDs(ids)
	}
	if request.Filters.AskID != nil {
		q = q.FilterByAtomicSwapAskID(*request.Filters.AskID)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateAtomicSwapBidRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest,
) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset) {
		bid := record.Details.CreateAtomicSwapBid
		asset, err := h.QuoteAssetsQ.FilterByIDs([]int64{int64(bid.AskID)}).
			FilterByCodes([]string{bid.QuoteAsset}).Get()
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get quote asset")
		}
		if asset == nil {
			return regources.ReviewableRequest{}, errors.New("quote asset not found")
		}

		resource := resources.NewAtomicSwapAskQuoteAsset(*asset)
		included.Add(&resource)
	}

	return resource, nil
}
