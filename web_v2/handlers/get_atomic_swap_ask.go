package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/history2"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources/generated"
)

// GetAtomicSwapAsk - processes request to get atomic swap bid and it's details by id
func GetAtomicSwapAsk(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAtomicSwapAskHandler{
		AssetsQ:               core2.NewAssetsQ(coreRepo),
		AtomicSwapBidQ:        core2.NewAtomicSwapAskQ(coreRepo),
		AtomicSwapQuoteAssetQ: core2.NewAtomicSwapQuoteAssetQ(coreRepo),
		BalanceQ:              history2.NewBalancesQ(ctx.HistoryRepo(r)),
		Log:                   ctx.Log(r),
	}

	request, err := requests.NewGetAtomicSwapAsk(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAtomicSwapAsk(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get atomic swap ask", logan.F{
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

type getAtomicSwapAskHandler struct {
	AssetsQ               core2.AssetsQ
	BalanceQ              history2.BalancesQ
	AtomicSwapBidQ        core2.AtomicSwapAskQ
	AtomicSwapQuoteAssetQ core2.AtomicSwapQuoteAssetQ
	Log                   *logan.Entry
}

// GetAtomicSwapAsk returns atomic swap bid with related resources
func (h *getAtomicSwapAskHandler) GetAtomicSwapAsk(request *requests.GetAtomicSwapAsk,
) (*regources.AtomicSwapAskResponse, error) {
	bid, err := h.AtomicSwapBidQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get atomic swap ask by id")
	}
	if bid == nil {
		return nil, nil
	}

	response := &regources.AtomicSwapAskResponse{
		Data: resources.NewAtomicSwapAsk(*bid),
	}

	bidOwner := resources.NewAccountKey(bid.OwnerID)
	response.Data.Relationships.Owner = bidOwner.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeAskOwner) {
		response.Included.Add(&bidOwner)
	}

	baseBalanceKey := resources.NewBalanceKey(bid.BaseBalanceID)
	response.Data.Relationships.BaseBalance = baseBalanceKey.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeAskBaseBalance) {
		baseBalance := regources.Balance{
			Key: baseBalanceKey,
			Relationships: &regources.BalanceRelationships{
				Asset: resources.NewAssetKey(bid.BaseAsset).AsRelation(),
			},
		}
		response.Included.Add(&baseBalance)
	}

	assetRaw, err := h.AssetsQ.GetByCode(bid.BaseAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get base asset")
	}
	if assetRaw == nil {
		return nil, errors.From(errors.New("base asset not found"), logan.F{
			"asset code": bid.BaseAsset,
		})
	}
	response.Data.Relationships.BaseAsset = resources.NewAssetKey(assetRaw.Code).AsRelation()
	if request.ShouldInclude(requests.IncludeTypeAskBaseAsset) {
		asset := resources.NewAsset(*assetRaw)
		response.Included.Add(&asset)
	}

	quoteAssetsRaw, err := h.AtomicSwapQuoteAssetQ.FilterByIDs([]int64{bid.AskID}).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select ask quote assets")
	}
	if quoteAssetsRaw == nil {
		return nil, errors.From(errors.New("expected ask quote assets to exists"), logan.F{
			"ask_id": bid.AskID,
		})
	}

	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(quoteAssetsRaw)),
	}

	for _, quoteAssetRaw := range quoteAssetsRaw {
		quoteAsset := resources.NewAtomicSwapAskQuoteAsset(quoteAssetRaw)
		quoteAssets.Data = append(quoteAssets.Data, quoteAsset.Key)

		if request.ShouldInclude(requests.IncludeTypeAskQuoteAssets) {
			response.Included.Add(&quoteAsset)
		}
	}
	response.Data.Relationships.QuoteAssets = quoteAssets

	return response, nil
}
