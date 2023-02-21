package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func GetDataCreationRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetDataCreationRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getDataCreationRequestsHandler{
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

type getDataCreationRequestsHandler struct {
	R         requests.GetDataCreationRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	Log       *logan.Entry
}

func (h *getDataCreationRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetDataCreationRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeDataCreation))

	if request.Filters.SecurityType != nil {
		q = h.RequestsQ.FilterByDataCreationSecurityType(*request.Filters.SecurityType)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getDataCreationRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	return resource, nil
}
