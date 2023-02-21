package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"

	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func GetCreatePollRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreatePollRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreatePollRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	constraints := []*string{
		request.GetRequestsBase.Filters.Requestor,
		request.GetRequestsBase.Filters.Reviewer,
	}

	if !isAllowed(r, w, constraints...) {
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

type getCreatePollRequestsHandler struct {
	R         requests.GetCreatePollRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

func (h *getCreatePollRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreatePollRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreatePoll))

	if request.Filters.PermissionType != nil {
		q = q.FilterByCreatePollPermissionType(*request.Filters.PermissionType)
	}

	if request.Filters.VoteConfirmationRequired != nil {
		q = q.FilterByCreatePollVoteConfirmationRequired(*request.Filters.VoteConfirmationRequired)
	}
	if request.Filters.ResultProvider != nil {
		q = q.FilterByCreatePollResultProvider(*request.Filters.ResultProvider)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreatePollRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	return resource, nil
}
