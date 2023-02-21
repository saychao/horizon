package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func GetRedemptionRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRedemptionRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getRedemptionRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		BalancesQ: core2.NewBalancesQ(coreRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
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

type getRedemptionRequestsHandler struct {
	R         requests.GetRedemptionRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	BalancesQ core2.BalancesQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

func (h *getRedemptionRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetRedemptionRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypePerformRedemption))

	if request.Filters.DestinationAccount != nil {
		q = q.FilterByAssetUpdateAsset(*request.Filters.DestinationAccount)
	}
	if request.Filters.SourceBalance != nil {
		q = q.FilterByAssetUpdateAsset(*request.Filters.SourceBalance)
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getRedemptionRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeRedemptionRequestsSourceBalance) {
		balance, err := h.BalancesQ.GetByAddress(record.Details.Redemption.SourceBalanceID)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get balance")
		}

		if balance == nil {
			return regources.ReviewableRequest{}, errors.New("balance not found")
		}
		resource := &regources.Balance{
			Key: resources.NewBalanceKey(balance.BalanceAddress),
		}
		included.Add(resource)
	}

	if h.R.ShouldInclude(requests.IncludeTypeRedemptionRequestsDestinationAccount) {
		account, err := h.AccountsQ.GetByAddress(record.Details.Redemption.DestinationAccountID)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get account")
		}

		if account == nil {
			return regources.ReviewableRequest{}, errors.New("account not found")
		}
		resource := &regources.Balance{
			Key: resources.NewBalanceKey(account.Address),
		}
		included.Add(resource)
	}

	return resource, nil
}
