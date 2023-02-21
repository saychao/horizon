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

// GetAssetPair - processes request to get asset pair and it's details by it's ID (`base:quote`)
func GetAssetPair(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAssetPairHandler{
		AssetPairsQ: core2.NewAssetPairsQ(coreRepo),
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetAssetPair(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAssetPair(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset pair", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getAssetPairHandler struct {
	AssetPairsQ core2.AssetPairsQ
	AssetsQ     core2.AssetsQ
	Log         *logan.Entry
}

// GetAssetPair returns asset pair with related resources
func (h *getAssetPairHandler) GetAssetPair(request *requests.GetAssetPair) (*regources.AssetPairResponse, error) {
	q := h.AssetPairsQ

	if request.ShouldInclude(requests.IncludeTypeAssetPairBaseAsset) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeAssetPairQuoteAsset) {
		q = q.WithQuoteAsset()
	}

	coreAssetPair, err := q.GetByBaseAndQuote(request.BaseAsset, request.QuoteAsset)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset pair by base and quote")
	}
	if coreAssetPair == nil {
		return nil, nil
	}

	response := &regources.AssetPairResponse{
		Data: resources.NewAssetPair(*coreAssetPair),
	}

	baseAssetKey := resources.NewAssetKey(coreAssetPair.Base)
	quoteAssetKey := resources.NewAssetKey(coreAssetPair.Quote)
	response.Data.Relationships.BaseAsset = baseAssetKey.AsRelation()
	response.Data.Relationships.QuoteAsset = quoteAssetKey.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeAssetPairBaseAsset) {
		coreBaseAsset := coreAssetPair.BaseAsset
		baseAsset := resources.NewAsset(*coreBaseAsset)

		response.Included.Add(&baseAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeAssetPairQuoteAsset) {
		coreQuoteAsset := coreAssetPair.QuoteAsset
		quoteAsset := resources.NewAsset(*coreQuoteAsset)

		response.Included.Add(&quoteAsset)
	}

	return response, nil
}
