package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/logan/v3"
	regources "gitlab.com/tokend/regources/generated"

	"github.com/saychao/horizon/db2/history2"

	"github.com/saychao/horizon/web_v2/ctx"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/saychao/horizon/web_v2/requests"
)

func GetDeferredPayment(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetDeferredPayment(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	hrepo := ctx.HistoryRepo(r)
	h := getDeferredPaymentHandler{
		DeferredPaymentQ: history2.NewDeferredPaymentQ(hrepo),
		AccountQ:         core2.NewAccountsQ(ctx.CoreRepo(r)),
		BalanceQ:         core2.NewBalancesQ(ctx.CoreRepo(r)),
		log:              ctx.Log(r),
	}

	result, err := h.DeferredPaymentQ.GetByID(request.DeferredPaymentID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get deferredPayment")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	source, err := h.BalanceQ.GetByAddress(result.SourceBalance)
	if err != nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, &source.AccountAddress, &result.DestinationAccount) {
		return
	}

	response := regources.DeferredPaymentResponse{
		Data: resources.NewDeferredPayment(*result),
	}

	ape.Render(w, response)
}

type getDeferredPaymentHandler struct {
	DeferredPaymentQ history2.DeferredPaymentQ
	AccountQ         core2.AccountsQ
	BalanceQ         core2.BalancesQ

	log *logan.Entry
}
