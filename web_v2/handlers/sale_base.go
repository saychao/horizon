package handlers

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

var (
	unexpectedParticipationsQType = errors.New("unexpected participationsQ type")
	unexpectedSaleType            = errors.New("unexpected sale type")
	unexpectedSaleState           = errors.New("unexpected sale state")
)

type getSaleBase struct {
	SalesQ           history2.SalesQ
	AssetsQ          history2.AssetQ
	saleCapConverter *saleCapConverter
	Log              *logan.Entry

	ParticipationQ history2.SaleParticipationQ
	OffersQ        core2.OffersQ
}

func (h *getSaleBase) getAndPopulateResponse(q history2.SalesQ, request *requests.GetSale) (*regources.SaleResponse, error) {
	historySale, err := q.Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sale by ID")
	}
	if historySale == nil {
		return nil, nil
	}

	err = h.saleCapConverter.PopulateSaleCap(historySale)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate sale cap")
	}

	response := &regources.SaleResponse{
		Data: resources.NewSale(*historySale),
	}

	defaultQuoteAsset := resources.NewSaleDefaultQuoteAsset(*historySale)
	response.Data.Relationships.DefaultQuoteAsset = defaultQuoteAsset.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeSaleDefaultQuoteAsset) {
		response.Included.Add(&defaultQuoteAsset)
	}

	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(historySale.QuoteAssets.QuoteAssets)),
	}

	for _, historyQuoteAsset := range historySale.QuoteAssets.QuoteAssets {
		quoteAsset := resources.NewSaleQuoteAsset(historyQuoteAsset, historySale.ID)
		quoteAssets.Data = append(quoteAssets.Data, quoteAsset.Key)

		if request.ShouldInclude(requests.IncludeTypeSaleQuoteAssets) {
			response.Included.Add(&quoteAsset)
		}
	}
	response.Data.Relationships.QuoteAssets = quoteAssets

	if request.ShouldInclude(requests.IncludeTypeSaleBaseAsset) {
		asset := resources.NewAssetV2(*historySale.Asset)
		response.Included.Add(&asset)
	}

	saleParticipationsCountMap, saleParticipantsCountMap, err := salesParticipationsInfo(h.ParticipationQ, h.OffersQ, *historySale)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load participations info")
	}

	if count, ok := saleParticipationsCountMap[historySale.ID]; ok {
		response.Data.Attributes.ParticipationsCount = count
	}

	if count, ok := saleParticipantsCountMap[historySale.ID]; ok {
		response.Data.Attributes.ParticipantsCount = count
	}

	return response, nil
}

type salesBaseHandler struct {
	SalesQ           history2.SalesQ
	AssetsQ          history2.AssetQ
	saleCapConverter *saleCapConverter
	Log              *logan.Entry

	ParticipationQ history2.SaleParticipationQ
	OffersQ        core2.OffersQ
}

func (h *salesBaseHandler) populateResponse(historySales []history2.Sale,
	request requests.SalesBase,
	response *regources.SaleListResponse) error {

	participationsCount, participantsCount, err := salesParticipationsInfo(h.ParticipationQ, h.OffersQ, historySales...)
	if err != nil {
		return errors.Wrap(err, "failed to populate participations count")
	}

	for _, historySale := range historySales {
		sale := resources.NewSale(historySale)

		err := h.saleCapConverter.PopulateSaleCap(&historySale)
		if err != nil {
			return errors.Wrap(err, "failed to populate sale cap")
		}
		sale.Relationships.QuoteAssets = &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(historySale.QuoteAssets.QuoteAssets)),
		}

		defaultQuoteAsset := resources.NewSaleDefaultQuoteAsset(historySale)
		sale.Relationships.DefaultQuoteAsset = defaultQuoteAsset.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeSaleDefaultQuoteAsset) {
			response.Included.Add(&defaultQuoteAsset)
		}

		for _, historyQuoteAsset := range historySale.QuoteAssets.QuoteAssets {
			quoteAsset := resources.NewSaleQuoteAsset(historyQuoteAsset, historySale.ID)
			sale.Relationships.QuoteAssets.Data = append(sale.Relationships.QuoteAssets.Data, quoteAsset.Key)

			if request.ShouldInclude(requests.IncludeTypeSaleQuoteAssets) {
				response.Included.Add(&quoteAsset)
			}
		}

		if request.ShouldInclude(requests.IncludeTypeSaleListBaseAssets) {
			asset := resources.NewAssetV2(*historySale.Asset)
			response.Included.Add(&asset)
		}

		if count, ok := participationsCount[historySale.ID]; ok {
			sale.Attributes.ParticipationsCount = count
		}

		if count, ok := participantsCount[historySale.ID]; ok {
			sale.Attributes.ParticipantsCount = count
		}

		response.Data = append(response.Data, sale)
	}

	return nil
}

func applySaleIncludes(s requests.SalesBase, q history2.SalesQ) history2.SalesQ {
	if s.ShouldInclude(requests.IncludeTypeSaleListBaseAssets) &&
		!s.ShouldFilter(requests.FilterTypeSaleListAssetType) {
		q = q.WithAsset()
	}

	return q
}

func applySaleFilters(s requests.SalesBase, q history2.SalesQ) history2.SalesQ {
	if s.ShouldFilter(requests.FilterTypeSaleListOwner) {
		q = q.FilterByOwner(s.Filters.Owner)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListBaseAsset) {
		q = q.FilterByBaseAsset(s.Filters.BaseAsset)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListAssetType) {
		q = q.FilterByAssetType(s.Filters.AssetType)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxEndTime) {
		q = q.FilterByMaxEndTime(*s.Filters.MaxEndTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxStartTime) {
		q = q.FilterByMaxStartTime(*s.Filters.MaxStartTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinStartTime) {
		q = q.FilterByMinStartTime(*s.Filters.MinStartTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinEndTime) {
		q = q.FilterByMinEndTime(*s.Filters.MinEndTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListState) {
		q = q.FilterByState(s.Filters.State)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListSaleType) {
		q = q.FilterBySaleType(s.Filters.SaleType)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinHardCap) {
		q = q.FilterByMinHardCap(uint64(s.Filters.MinHardCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinSoftCap) {
		q = q.FilterByMinSoftCap(uint64(s.Filters.MinSoftCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxHardCap) {
		q = q.FilterByMaxHardCap(uint64(s.Filters.MaxHardCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxSoftCap) {
		q = q.FilterByMaxSoftCap(uint64(s.Filters.MaxSoftCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListIDs) {
		q = q.FilterByIDs(s.Filters.IDs)
	}

	if s.ShouldFilter(requests.FilterTypeSaleCompanyID) {
		q = q.FilterByCompanyIDs(s.Filters.CompanyIDs)
	}

	return q
}

type PrtQType int

const (
	closedPrtQ PrtQType = iota
	pendingPrtQ
	undefinedPrtQ
	undefinedPrtQSaleType
	cancelledPrtQ
)

type SaleParticipationsInfo struct {
	SalesIDs   []uint64
	BaseAssets []string
	Owners     []string
}

func salesParticipationsInfo(saleParticipationsQ history2.SaleParticipationQ, offersQ core2.OffersQ, historySales ...history2.Sale) (core2.SaleParticipantsMap, core2.SaleParticipantsMap, error) {
	prtQTypeSalesMap := make(map[PrtQType]SaleParticipationsInfo)
	for _, sale := range historySales {
		prtQType := getParticipationsQType(sale)

		switch prtQType {
		case closedPrtQ:
			prtQTypeSalesMap[prtQType] = SaleParticipationsInfo{
				SalesIDs:   append(prtQTypeSalesMap[prtQType].SalesIDs, sale.ID),
				BaseAssets: append(prtQTypeSalesMap[prtQType].BaseAssets, sale.BaseAsset),
				Owners:     append(prtQTypeSalesMap[prtQType].Owners, sale.OwnerAddress),
			}
		case pendingPrtQ:
			prtQTypeSalesMap[prtQType] = SaleParticipationsInfo{
				SalesIDs: append(prtQTypeSalesMap[prtQType].SalesIDs, sale.ID),
			}
		case cancelledPrtQ:
			continue
		case undefinedPrtQSaleType:
			return nil, nil, unexpectedSaleType
		case undefinedPrtQ:
			return nil, nil, unexpectedSaleState
		default:
			return nil, nil, unexpectedParticipationsQType
		}
	}

	participationsCount := make([]core2.SaleIDParticipantsCount, 0, len(historySales))
	participantsCount := make([]core2.SaleIDParticipantsCount, 0, len(historySales))
	prtQSalesMap := make(map[PrtQType]participationsQ)
	for prtType, salesInfo := range prtQTypeSalesMap {
		if _, ok := prtQSalesMap[prtType]; !ok {
			switch prtType {
			case pendingPrtQ:
				prtQSalesMap[prtType] = pendingParticipationsQ{
					offersQ: offersQ.
						FilterByIsBuy(true).
						FilterByOrderBookIDs(salesInfo.SalesIDs...),
				}
			case closedPrtQ:
				prtQSalesMap[prtType] = closedParticipationsQ{
					participationQ: saleParticipationsQ.
						FilterBySaleIDs(salesInfo.SalesIDs...).
						FilterBySaleBaseAssets(salesInfo.BaseAssets...).
						FilterByNotSaleOwners(salesInfo.Owners...),
				}
			default:
				return nil, nil, unexpectedParticipationsQType
			}
		}

		participationsCountTemp, err := prtQSalesMap[prtType].ParticipationsCount()
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to select participations count")
		}
		participantsCountTemp, err := prtQSalesMap[prtType].ParticipantsCount()
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to select participants count")
		}

		participationsCount = append(participationsCount, participationsCountTemp...)
		participantsCount = append(participantsCount, participantsCountTemp...)
	}

	return core2.SaleParticipantsToMap(participationsCount), core2.SaleParticipantsToMap(participantsCount), nil
}

func getParticipationsQType(historySale history2.Sale) PrtQType {
	switch historySale.State {
	case regources.SaleStateCanceled:
		return cancelledPrtQ
	case regources.SaleStateOpen:
		switch historySale.SaleType {
		case xdr.SaleTypeImmediate:
			return closedPrtQ
		case xdr.SaleTypeBasicSale, xdr.SaleTypeCrowdFunding, xdr.SaleTypeFixedPrice:
			return pendingPrtQ
		default:
			return undefinedPrtQSaleType
		}
	case regources.SaleStateClosed:
		return closedPrtQ
	default:
		return undefinedPrtQ
	}
}
