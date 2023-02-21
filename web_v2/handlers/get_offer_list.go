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

// GetOfferList - processes request to get the list of offers
func GetOfferList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getOfferListHandler{
		OffersQ:   core2.NewOffersQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetOfferList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, request.Filters.Owner) {
		return
	}

	result, err := handler.GetOfferList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get offer list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getOfferListHandler struct {
	OffersQ   core2.OffersQ
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetOfferList returns offer with related resources
func (h *getOfferListHandler) GetOfferList(request *requests.GetOfferList) (*regources.OfferListResponse, error) {
	q := h.OffersQ.Page(request.PageParams)

	if request.Filters.BaseBalance != nil {
		q = q.FilterByBaseBalanceID(*request.Filters.BaseBalance)
	}

	if request.Filters.QuoteBalance != nil {
		q = q.FilterByQuoteBalanceID(*request.Filters.QuoteBalance)
	}

	if request.Filters.BaseAsset != nil {
		q = q.FilterByBaseAssetCode(*request.Filters.BaseAsset)
	}

	if request.Filters.QuoteAsset != nil {
		q = q.FilterByQuoteAssetCode(*request.Filters.QuoteAsset)
	}

	if request.Filters.Owner != nil {
		q = q.FilterByOwnerID(*request.Filters.Owner)
	}

	if request.Filters.OrderBook != nil {
		q = q.FilterByOrderBookID(*request.Filters.OrderBook)
	}

	if request.Filters.IsBuy != nil {
		q = q.FilterByIsBuy(*request.Filters.IsBuy)
	}

	if request.ShouldInclude(requests.IncludeTypeOfferListBaseAssets) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOfferListQuoteAssets) {
		q = q.WithQuoteAsset()
	}

	coreOffers, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get offer list")
	}

	response := &regources.OfferListResponse{
		Data:  make([]regources.Offer, 0, len(coreOffers)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, coreOffer := range coreOffers {
		offer := resources.NewOffer(coreOffer)

		offer.Relationships.Owner = resources.NewAccountKey(coreOffer.OwnerID).AsRelation()
		offer.Relationships.BaseAsset = resources.NewAssetKey(coreOffer.BaseAssetCode).AsRelation()
		offer.Relationships.QuoteAsset = resources.NewAssetKey(coreOffer.QuoteAssetCode).AsRelation()
		offer.Relationships.BaseBalance = resources.NewBalanceKey(coreOffer.BaseBalanceID).AsRelation()
		offer.Relationships.QuoteBalance = resources.NewBalanceKey(coreOffer.QuoteBalanceID).AsRelation()

		response.Data = append(response.Data, offer)

		if request.ShouldInclude(requests.IncludeTypeOfferListBaseAssets) {
			baseAsset := resources.NewAsset(*coreOffer.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOfferListQuoteAssets) {
			quoteAsset := resources.NewAsset(*coreOffer.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	return response, nil
}
