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

func GetCreateAtomicSwapAskRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateAtomicSwapAskRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateAtomicSwapAskRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		BalancesQ: history2.NewBalancesQ(historyRepo),
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

type getCreateAtomicSwapAskRequestsHandler struct {
	R         requests.GetCreateAtomicSwapAskRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	BalancesQ history2.BalancesQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateAtomicSwapAskRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateAtomicSwapAskRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateAtomicSwapAsk))

	if request.Filters.BaseBalance != nil {
		q = q.FilterByCreateAtomicSwapAskBalance(*request.Filters.BaseBalance)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateAtomicSwapAskRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest,
) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateAtomicSwapAskRequestsBalance) {
		balance, err := h.BalancesQ.GetByAddress(record.Details.CreateAtomicSwapAsk.BaseBalance)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get balance")
		}

		if balance == nil {
			return regources.ReviewableRequest{}, errors.New("balance not found")
		}
		resource := &regources.Balance{
			Key: resources.NewBalanceKey(balance.Address),
		}
		included.Add(resource)
	}

	if h.R.ShouldInclude(requests.IncludeTypeCreateAtomicSwapAskRequestsQuoteAssets) {
		for _, record := range record.Details.CreateAtomicSwapAsk.QuoteAssets {
			quoteAsset := core2.AtomicSwapQuoteAsset{
				AskID:      0,
				QuoteAsset: record.Asset,
				Price:      uint64(record.Price),
			}
			asset := resources.NewAtomicSwapAskQuoteAsset(quoteAsset)
			included.Add(&asset)
		}
	}

	return resource, nil
}
