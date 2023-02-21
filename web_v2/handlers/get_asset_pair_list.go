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

// GetAssetPairList - processes request to get the list of asset pairs and their details
func GetAssetPairList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetPairListHandler{
		AssetPairsQ: core2.NewAssetPairsQ(coreRepo),
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetAssetPairList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetPairList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAssetPairListHandler struct {
	AssetPairsQ core2.AssetPairsQ
	AssetsQ     core2.AssetsQ
	Log         *logan.Entry
}

// GetAssetPairList returns asset pair list with related resources
func (h *getAssetPairListHandler) GetAssetPairList(request *requests.GetAssetPairList) (*regources.AssetPairListResponse, error) {
	q := h.AssetPairsQ.Page(request.PageParams)

	if request.Filters.Asset != nil {
		q = q.FilterByAsset(*request.Filters.Asset)
	}

	if request.Filters.BaseAsset != nil {
		q = q.FilterByBaseAsset(*request.Filters.BaseAsset)
	}

	if request.Filters.QuoteAsset != nil {
		q = q.FilterByQuoteAsset(*request.Filters.QuoteAsset)
	}

	if request.Filters.Policy != nil {
		q = q.FilterByPolicy(*request.Filters.Policy)
	}

	if request.ShouldInclude(requests.IncludeTypeAssetPairListBaseAssets) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeAssetPairListQuoteAssets) {
		q = q.WithQuoteAsset()
	}

	coreAssetPairs, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset pair list")
	}

	response := &regources.AssetPairListResponse{
		Data:  make([]regources.AssetPair, 0, len(coreAssetPairs)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for i := range coreAssetPairs {
		assetPair := resources.NewAssetPair(coreAssetPairs[i])

		baseAssetKey := resources.NewAssetKey(coreAssetPairs[i].Base)
		quoteAssetKey := resources.NewAssetKey(coreAssetPairs[i].Quote)
		assetPair.Relationships.BaseAsset = baseAssetKey.AsRelation()
		assetPair.Relationships.QuoteAsset = quoteAssetKey.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeAssetPairListBaseAssets) {
			coreBaseAsset := coreAssetPairs[i].BaseAsset
			baseAsset := resources.NewAsset(*coreBaseAsset)

			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeAssetPairListQuoteAssets) {
			coreQuoteAsset := coreAssetPairs[i].QuoteAsset
			quoteAsset := resources.NewAsset(*coreQuoteAsset)

			response.Included.Add(&quoteAsset)
		}

		response.Data = append(response.Data, assetPair)
	}

	return response, nil
}
