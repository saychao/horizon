package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetAssetList - processes request to get the list of assets
func GetAssetList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetListHandler{
		AccountsQ: core2.NewAccountsQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetAssetList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAssetListHandler struct {
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetAssetList returns the list of assets with related resources
func (h *getAssetListHandler) GetAssetList(request *requests.GetAssetList) (*regources.AssetListResponse, error) {
	q := h.AssetsQ.Page(request.PageParams)
	if request.Owner != nil {
		q = q.FilterByOwner(*request.Owner)
	}
	if request.Policy != nil {
		q = q.FilterByPolicy(*request.Policy)
	}
	if request.State != nil {
		q = q.FilterByState(*request.State)
	}
	if request.Codes != nil {
		q = q.FilterByCodes(request.Codes)
	}
	if request.Types != nil {
		q = q.FilterByTypes(request.Types)
	}

	assets, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset list")
	}

	response := &regources.AssetListResponse{
		Data:  make([]regources.Asset, 0, len(assets)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for i := range assets {
		asset := resources.NewAsset(assets[i])
		owner := resources.NewAccountKey(assets[i].Owner)
		asset.Relationships.Owner = owner.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeAssetListOwners) {
			response.Included.Add(&owner)
		}

		response.Data = append(response.Data, asset)
	}

	return response, nil
}
