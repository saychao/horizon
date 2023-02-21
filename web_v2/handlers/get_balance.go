package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)

	handler := getBalanceHandler{
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		BalancesQ: core2.NewBalancesQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetBalance(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetBalance(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get balance", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if request.ShouldInclude(
		requests.IncludeTypeBalanceState,
	) {
		if !isAllowed(r, w, &result.Data.Relationships.Owner.Data.ID) {
			return
		}
	}

	ape.Render(w, result)
}

type getBalanceHandler struct {
	BalancesQ core2.BalancesQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

func (h *getBalanceHandler) GetBalance(request *requests.GetBalance) (*regources.BalanceResponse, error) {
	balance, err := h.BalancesQ.GetByAddress(request.BalanceID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot filter balances by ID", logan.F{
			"id": request.BalanceID,
		})
	}
	if balance == nil {
		return nil, nil
	}

	response := regources.BalanceResponse{
		Data: *resources.NewBalance(balance),
	}

	response.Data.Relationships = &regources.BalanceRelationships{
		Owner: resources.NewAccountKey(balance.AccountAddress).AsRelation(),
		Asset: resources.NewAssetKey(balance.AssetCode).AsRelation(),
		State: resources.NewBalanceState(balance).AsRelation(),
	}

	if request.ShouldInclude(requests.IncludeTypeBalanceAsset) {
		if err = h.includeAsset(balance.AssetCode, &response.Included); err != nil {
			return nil, errors.Wrap(err, "failed to include asset")
		}
	}

	if request.ShouldInclude(requests.IncludeTypeBalanceState) {
		response.Included.Add(resources.NewBalanceState(balance))
	}

	return &response, nil
}

func (h *getBalanceHandler) includeAsset(assetCode string, included *regources.Included) error {
	asset, err := h.AssetsQ.GetByCode(assetCode)
	if err != nil {
		return errors.Wrap(err, "cannot get asset by code", logan.F{
			"asset_code": assetCode,
		})
	}

	assetResource := resources.NewAsset(*asset)
	included.Add(&assetResource)

	return nil
}
