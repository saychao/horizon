package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func GetChangeRoleRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetChangeRoleRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getChangeRoleRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	if !isAllowed(r, w, request.GetRequestsBase.Filters.Requestor, request.GetRequestsBase.Filters.Reviewer, request.Filters.Account) {
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

type getChangeRoleRequestsHandler struct {
	R         requests.GetChangeRoleRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

func (h *getChangeRoleRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetChangeRoleRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeChangeRole))

	if acc := request.Filters.Account; acc != nil {
		q = q.FilterByChangeRoleAccount(*acc)
	}
	if role := request.Filters.AccountRole; role != nil {
		q = q.FilterByChangeRoleToSet(*role)
	}
	if reg := request.Filters.AdminRegistered; reg != nil {
		q = q.FilterByRegistrationType(*reg)
	}
	if ts := request.Filters.TimestampForZoho; ts != nil {
		q = q.FilterByTimestampGT(*ts)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getChangeRoleRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeChangeRoleRequestsAccount) {
		account, err := h.AccountsQ.GetByAddress(record.Details.ChangeRole.DestinationAccount)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get account")
		}

		if account == nil {
			return regources.ReviewableRequest{}, errors.New("account not found")
		}
		resource := resources.NewAccount(*account, nil)
		included.Add(&resource)
	}

	return resource, nil
}
