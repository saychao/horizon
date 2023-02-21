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

// GetOffer - processes request to get offer and it's details by it's ID
func GetOffer(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getOfferHandler{
		Log:       ctx.Log(r),
		OffersQ:   core2.NewOffersQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
	}

	request, err := requests.NewGetOffer(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	coreOffer, err := handler.GetOffer(*request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get offer", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if coreOffer == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, &coreOffer.OwnerID) {
		return
	}

	handler.RenderResponse(r, w, request, coreOffer)
}

type getOfferHandler struct {
	OffersQ   core2.OffersQ
	AssetsQ   core2.AssetsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

// GetOffer returns gets offer from database
func (h *getOfferHandler) GetOffer(request requests.GetOffer) (*core2.Offer, error) {
	q := h.OffersQ

	if request.ShouldInclude(requests.IncludeTypeOfferBaseAsset) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOfferQuoteAsset) {
		q = q.WithQuoteAsset()
	}

	coreOffer, err := q.GetByOfferID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get offer by ID")
	}

	return coreOffer, nil
}

// RenderResponse renders offer with related resources
func (h *getOfferHandler) RenderResponse(r *http.Request, w http.ResponseWriter,
	request *requests.GetOffer, coreOffer *core2.Offer) {

	response := regources.OfferResponse{
		Data: resources.NewOffer(*coreOffer),
	}

	response.Data.Relationships.Owner = resources.NewAccountKey(coreOffer.OwnerID).AsRelation()
	response.Data.Relationships.BaseAsset = resources.NewAssetKey(coreOffer.BaseAssetCode).AsRelation()
	response.Data.Relationships.QuoteAsset = resources.NewAssetKey(coreOffer.QuoteAssetCode).AsRelation()
	response.Data.Relationships.BaseBalance = resources.NewBalanceKey(coreOffer.BaseBalanceID).AsRelation()
	response.Data.Relationships.QuoteBalance = resources.NewBalanceKey(coreOffer.QuoteBalanceID).AsRelation()

	if request.ShouldInclude(requests.IncludeTypeOfferBaseAsset) {
		baseAsset := resources.NewAsset(*coreOffer.BaseAsset)
		response.Included.Add(&baseAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeOfferQuoteAsset) {
		quoteAsset := resources.NewAsset(*coreOffer.QuoteAsset)
		response.Included.Add(&quoteAsset)
	}

	ape.Render(w, response)
}
