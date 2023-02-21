package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

// GetSaleParticipations - processes request to get list of sale participations
func GetSaleParticipations(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetSaleParticipations(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getSaleParticipationsHandler{
		ParticipationQ: history2.NewSaleParticipationQ(ctx.HistoryRepo(r)),
		OffersQ:        core2.NewOffersQ(ctx.CoreRepo(r)),
		AssetsQ:        core2.NewAssetsQ(ctx.CoreRepo(r)),
		SalesQ:         history2.NewSalesQ(ctx.HistoryRepo(r)),
		Log:            ctx.Log(r),
	}

	sale, err := handler.getSale(request.SaleID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if sale == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, &sale.OwnerAddress) {
		return
	}

	result, err := handler.GetSaleParticipations(sale, request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale participations", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSaleParticipationsHandler struct {
	Log            *logan.Entry
	SalesQ         history2.SalesQ
	OffersQ        core2.OffersQ
	AssetsQ        core2.AssetsQ
	ParticipationQ history2.SaleParticipationQ
}

type participationsQ interface {
	// FilterByParticipant - filters out participations by participant address
	FilterByParticipant(id string) participationsQ
	// FilterByQuoteAsset - filters out participations by quote asset
	FilterByQuoteAsset(code string) participationsQ
	// Select - select records from db and wraps them to participations
	Select() ([]regources.SaleParticipation, error)
	// ParticipationsCount - get slice of sales ids mapped to participations count
	ParticipationsCount() ([]core2.SaleIDParticipantsCount, error)
	// ParticipantsCount - get slice of sales ids mapped to participants count
	ParticipantsCount() ([]core2.SaleIDParticipantsCount, error)
}

// GetSaleParticipations returns sale with related resources
func (h *getSaleParticipationsHandler) GetSaleParticipations(sale *history2.Sale, request *requests.GetSaleParticipations,
) (*regources.SaleParticipationListResponse, error) {
	response := regources.SaleParticipationListResponse{
		Data: make([]regources.SaleParticipation, 0),
	}

	var q participationsQ

	switch sale.State {
	case regources.SaleStateCanceled:
		return &response, nil
	case regources.SaleStateOpen:
		switch sale.SaleType {
		case xdr.SaleTypeImmediate:
			// on immediate sale offers matched right away after creating participation, so we can use only history
			q = newClosedParticipationQ(request, h.ParticipationQ, sale)
		case xdr.SaleTypeBasicSale, xdr.SaleTypeCrowdFunding, xdr.SaleTypeFixedPrice:
			q = newPendingParticipationQ(request, h.OffersQ)
		default:
			return nil, errors.From(errors.New("unexpected sale type"), logan.F{
				"sale_type": sale.SaleType.String(),
			})
		}
	case regources.SaleStateClosed:
		q = newClosedParticipationQ(request, h.ParticipationQ, sale)
	default:
		return nil, errors.From(errors.New("unexpected sale state"), logan.F{
			"sale_state": sale.State,
		})
	}

	if request.Filters.Participant != nil {
		q = q.FilterByParticipant(*request.Filters.Participant)
	}

	if request.Filters.QuoteAsset != nil {
		q = q.FilterByQuoteAsset(*request.Filters.QuoteAsset)
	}

	var err error
	response.Data, err = q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load participations")
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(request.PageParams, "")
	}

	if len(response.Data) == 0 {
		return &response, nil
	}

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationsBaseAsset) {
		baseAsset, err := h.getAsset(sale.BaseAsset)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get asset")
		}

		response.Included.Add(baseAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeSaleParticipationsQuoteAsset) {
		assets, err := h.getAssets(response.Data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get assets")
		}

		for _, asset := range assets {
			response.Included.Add(&asset)
		}
	}

	return &response, nil
}

func (h *getSaleParticipationsHandler) getAssets(participations []regources.SaleParticipation) ([]regources.Asset, error) {
	codesSet := make(map[string]struct{})
	for _, participation := range participations {
		codesSet[participation.Relationships.QuoteAsset.Data.ID] = struct{}{}
	}

	codes := make([]string, 0)
	for code := range codesSet {
		codes = append(codes, code)
	}

	assets, err := h.AssetsQ.FilterByCodes(codes).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load assets")
	}

	if len(assets) != len(codes) {
		return nil, errors.From(errors.New("some assets have not been found"), logan.F{
			"asset_codes": codes,
		})
	}

	result := make([]regources.Asset, 0, len(assets))
	for _, asset := range assets {
		result = append(result, resources.NewAsset(asset))
	}

	return result, nil
}

func (h *getSaleParticipationsHandler) getAsset(code string) (*regources.Asset, error) {
	coreAsset, err := h.AssetsQ.GetByCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load asset by code", logan.F{
			"code": code,
		})
	}

	if coreAsset == nil {
		return nil, errors.From(errors.New("asset not found"), logan.F{
			"code": code,
		})
	}

	asset := resources.NewAsset(*coreAsset)

	return &asset, nil
}

func (h *getSaleParticipationsHandler) getSale(id uint64) (*history2.Sale, error) {
	sale, err := h.SalesQ.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale from db", logan.F{
			"id": id,
		})
	}

	if sale == nil {
		return nil, nil
	}

	return sale, nil
}
