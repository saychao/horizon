package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/saychao/horizon/db2/history2"

	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetSwap(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetSwap(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getSwapHandler{
		SwapsQ:    history2.NewSwapsQ(ctx.HistoryRepo(r)),
		AssetsQ:   core2.NewAssetsQ(ctx.CoreRepo(r)),
		BalancesQ: history2.NewBalancesQ(ctx.HistoryRepo(r)),
		Log:       ctx.Log(r),
	}

	result, err := handler.getSwap(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get swap", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if request.ShouldIncludeAny(
		requests.IncludeTypeSwapSourceBalance,
		requests.IncludeTypeSwapDestinationBalance,
	) {
		if !isAllowed(r, w, &result.Data.Relationships.Source.Data.ID, &result.Data.Relationships.Destination.Data.ID) {
			return
		}
	}

	ape.Render(w, result)
}

type getSwapHandler struct {
	SwapsQ    history2.SwapsQ
	AssetsQ   core2.AssetsQ
	BalancesQ history2.BalancesQ
	Log       *logan.Entry
}

// GetSale returns sale with related resources
func (h *getSwapHandler) getSwap(request *requests.GetSwap) (*regources.SwapResponse, error) {

	record, err := h.SwapsQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get swap")
	}

	if record == nil {
		return nil, nil
	}

	resource := resources.NewSwap(*record)
	response := &regources.SwapResponse{
		Data: resource,
	}
	if request.ShouldInclude(requests.IncludeTypeSwapAsset) {
		histAsset, err := h.AssetsQ.GetByCode(record.Asset)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get asset for swap")
		}
		if histAsset == nil {
			return nil, errors.New("Expected asset to exist")
		}
		asset := resources.NewAsset(*histAsset)
		response.Included.Add(&asset)
	}

	if request.ShouldInclude(requests.IncludeTypeSwapSourceBalance) {
		histBalance, err := h.BalancesQ.GetByAddress(record.SourceBalance)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get source balance for swap")
		}
		if histBalance == nil {
			return nil, errors.New("Expected balance to exist")
		}
		balance := &regources.Balance{
			Key: resources.NewBalanceKey(histBalance.Address),
		}
		response.Included.Add(balance)
	}

	if request.ShouldInclude(requests.IncludeTypeSwapDestinationBalance) {
		histBalance, err := h.BalancesQ.GetByAddress(record.DestinationAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get destination balance for swap")
		}
		if histBalance == nil {
			return nil, errors.New("Expected balance to exist")
		}
		balance := &regources.Balance{
			Key: resources.NewBalanceKey(histBalance.Address),
		}
		response.Included.Add(balance)
	}

	return response, nil
}
