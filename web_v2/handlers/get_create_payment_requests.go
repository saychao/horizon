package handlers

import (
	"net/http"

	"gitlab.com/tokend/go/xdr"

	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	regources "gitlab.com/tokend/regources/generated"
)

func GetCreatePaymentRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreatePaymentRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getCreatePaymentRequestsHandler{
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

type getCreatePaymentRequestsHandler struct {
	R         requests.GetCreatePaymentRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	Log       *logan.Entry
}

func (h *getCreatePaymentRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreatePaymentRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreatePayment))

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreatePaymentRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest,
) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	return resource, nil
}
