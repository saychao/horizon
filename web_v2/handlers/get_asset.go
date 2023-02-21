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

// GetAsset - processes request to get asset and it's details by asset code
func GetAsset(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetHandler{
		AccountsQ: core2.NewAccountsQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetAsset(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAsset(request)
	if err != nil {
		ctx.Log(r).WithError(err).WithField("request", request).Error("failed to get asset")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getAssetHandler struct {
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetAsset returns asset with related resources
func (h *getAssetHandler) GetAsset(request *requests.GetAsset) (*regources.AssetResponse, error) {
	asset, err := h.AssetsQ.GetByCode(request.Code)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset by code")
	}
	if asset == nil {
		return nil, nil
	}

	assetResponse := resources.NewAsset(*asset)
	response := &regources.AssetResponse{
		Data: assetResponse,
	}

	assetOwner := resources.NewAccountKey(asset.Owner)
	response.Data.Relationships.Owner = assetOwner.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeAssetOwner) {
		response.Included.Add(&assetOwner)
	}

	return response, nil
}
