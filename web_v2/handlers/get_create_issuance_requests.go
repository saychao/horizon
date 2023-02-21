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

func GetCreateIssuanceRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateIssuanceRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateIssuanceRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	constraints := []*string{
		request.GetRequestsBase.Filters.Requestor,
		request.GetRequestsBase.Filters.Reviewer,
	}

	// receiving balance owner should be able to see issuance requests
	if request.Filters.Receiver != nil && *request.Filters.Receiver != "" {
		balance, err := history2.NewBalancesQ(historyRepo).GetByAddress(*request.Filters.Receiver)
		if err != nil {
			ctx.Log(r).
				WithError(err).
				WithFields(logan.F{"receiver": *request.Filters.Receiver}).
				Error("failed to get receiver balance")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if balance == nil {
			ape.RenderErr(w, problems.NotFound())
			return
		}
		account, err := history2.NewAccountsQ(historyRepo).ByID(balance.AccountID)
		if err != nil {
			ctx.Log(r).
				WithError(err).
				WithFields(logan.F{"receiver": *request.Filters.Receiver}).
				Error("failed to get receiver account")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		constraints = append(constraints, &account.Address)

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

type getCreateIssuanceRequestsHandler struct {
	R         requests.GetCreateIssuanceRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getCreateIssuanceRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateIssuanceRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateIssuance))

	if len(request.Filters.Asset) != 0 {
		q = q.FilterByCreateIssuanceAsset(request.Filters.Asset...)
	}

	if request.Filters.Receiver != nil && *request.Filters.Receiver != "" {
		q = q.FilterByCreateIssuanceReceiver(*request.Filters.Receiver)
	}

	if request.Filters.Who != nil && *request.Filters.Who == "any" {
		q = q.FilterWhoNotNull()
	}

	return h.Base.SelectAndRender(w, request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateIssuanceRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateIssuanceRequestsAsset) {
		asset, err := h.AssetsQ.GetByCode(record.Details.CreateIssuance.Asset)
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
