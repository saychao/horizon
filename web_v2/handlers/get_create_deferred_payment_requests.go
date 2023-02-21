package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	regources "gitlab.com/tokend/regources/generated"
)

func GetCreateDeferredPaymentRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetDeferredPaymentCreateRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getCreateDeferredPaymentRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
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

type getCreateDeferredPaymentRequestsHandler struct {
	R         requests.GetDeferredPaymentCreateRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	Log       *logan.Entry
}

func (h *getCreateDeferredPaymentRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetDeferredPaymentCreateRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateDeferredPayment))

	if request.Filters.Destination != nil {
		q = h.RequestsQ.FilterByCreateDeferredPaymentDestination(*request.Filters.Destination)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateDeferredPaymentRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	return resource, nil
}
