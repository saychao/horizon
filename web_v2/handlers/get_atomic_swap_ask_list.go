package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/history2"

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

// GetAtomicSwapAskList - processes request to get the list of atomic swap bids
func GetAtomicSwapAskList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAtomicSwapAskListHandler{
		AtomicSwapAskQ:        core2.NewAtomicSwapAskQ(coreRepo),
		AtomicSwapQuoteAssetQ: core2.NewAtomicSwapQuoteAssetQ(coreRepo),
		BalanceQ:              history2.NewBalancesQ(ctx.HistoryRepo(r)),
		AssetsQ:               core2.NewAssetsQ(coreRepo),
		Log:                   ctx.Log(r),
	}

	request, err := requests.NewGetAtomicSwapAskList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAtomicSwapAskList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get atomic swap ask list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAtomicSwapAskListHandler struct {
	AssetsQ               core2.AssetsQ
	BalanceQ              history2.BalancesQ
	AtomicSwapAskQ        core2.AtomicSwapAskQ
	AtomicSwapQuoteAssetQ core2.AtomicSwapQuoteAssetQ
	Log                   *logan.Entry
}

// GetAtomicSwapAskList returns the list of atomic swap bids with related resources
func (h *getAtomicSwapAskListHandler) GetAtomicSwapAskList(request *requests.GetAtomicSwapAskList,
) (*regources.AtomicSwapAskListResponse, error) {
	q := h.AtomicSwapAskQ.Page(request.PageParams)
	if request.Filters.Owner != nil {
		q = q.FilterByOwner(*request.Filters.Owner)
	}
	if request.Filters.BaseAsset != nil {
		q = q.FilterByBaseAssets(request.Filters.BaseAsset)
	}
	if request.Filters.QuoteAssets != nil {
		q = q.FilterByQuoteAssets(request.Filters.QuoteAssets)
	}
	asks, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset list")
	}

	response := &regources.AtomicSwapAskListResponse{
		Data:  make([]regources.AtomicSwapAsk, 0, len(asks)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, ask := range asks {
		data := resources.NewAtomicSwapAsk(ask)
		bidOwner := resources.NewAccountKey(ask.OwnerID)
		data.Relationships.Owner = bidOwner.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeAskOwner) {
			response.Included.Add(&bidOwner)
		}

		baseBalanceKey := resources.NewBalanceKey(ask.BaseBalanceID)
		data.Relationships.BaseBalance = baseBalanceKey.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeAskBaseBalance) {
			baseBalance := regources.Balance{
				Key: baseBalanceKey,
				Relationships: &regources.BalanceRelationships{
					Asset: resources.NewAssetKey(ask.BaseAsset).AsRelation(),
				},
			}
			response.Included.Add(&baseBalance)
		}
		assetRaw, err := h.AssetsQ.GetByCode(ask.BaseAsset)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get base asset")
		}
		if assetRaw == nil {
			return nil, errors.From(errors.New("base asset not found"), logan.F{
				"asset code": ask.BaseAsset,
			})
		}
		data.Relationships.BaseAsset = resources.NewAssetKey(assetRaw.Code).AsRelation()
		if request.ShouldInclude(requests.IncludeTypeAskBaseAsset) {
			asset := resources.NewAsset(*assetRaw)
			response.Included.Add(&asset)
		}

		quoteAssetsRaw, err := h.AtomicSwapQuoteAssetQ.FilterByIDs([]int64{ask.AskID}).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to select ask quote assets")
		}
		if quoteAssetsRaw == nil {
			return nil, errors.From(errors.New("expected ask quote assets to exists"), logan.F{
				"bid_id": ask.AskID,
			})
		}

		quoteAssets := &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(quoteAssetsRaw)),
		}

		for _, quoteAssetRaw := range quoteAssetsRaw {
			quoteAsset := resources.NewAtomicSwapAskQuoteAsset(quoteAssetRaw)
			quoteAssets.Data = append(quoteAssets.Data, quoteAsset.Key)

			if request.ShouldInclude(requests.IncludeTypeAskListQuoteAssets) {
				response.Included.Add(&quoteAsset)
			}
		}
		data.Relationships.QuoteAssets = quoteAssets

		response.Data = append(response.Data, data)
	}

	return response, nil
}
