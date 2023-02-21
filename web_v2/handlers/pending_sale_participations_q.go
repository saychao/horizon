package handlers

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

type pendingParticipationsQ struct {
	offersQ core2.OffersQ
}

func newPendingParticipationQ(request *requests.GetSaleParticipations, q core2.OffersQ) pendingParticipationsQ {
	q = q.
		FilterByIsBuy(true).
		FilterByOrderBookID(int64(request.SaleID)).
		CursorPage(request.PageParams)

	return pendingParticipationsQ{
		offersQ: q,
	}
}

// FilterByParticipant - filters out participations by participant address
func (p pendingParticipationsQ) FilterByParticipant(address string) participationsQ {
	p.offersQ = p.offersQ.FilterByOwnerID(address)
	return p
}

// FilterByQuoteAsset - filters out participations by quote asset
func (p pendingParticipationsQ) FilterByQuoteAsset(code string) participationsQ {
	p.offersQ = p.offersQ.FilterByQuoteAssetCode(code)
	return p
}

// Select - select records from db and wraps them to participations
func (p pendingParticipationsQ) Select() ([]regources.SaleParticipation, error) {
	offers, err := p.offersQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load offers from db")
	}

	result := make([]regources.SaleParticipation, 0, len(offers))

	for _, offer := range offers {
		result = append(result, resources.NewSaleParticipation(
			offer.OfferID,
			offer.OwnerID,
			offer.BaseAssetCode,
			offer.QuoteAssetCode,
			offer.QuoteAmount),
		)
	}

	return result, nil
}

// ParticipationsCount - returns slice of sales ids mapped to participants count
func (p pendingParticipationsQ) ParticipationsCount() ([]core2.SaleIDParticipantsCount, error) {
	return p.offersQ.SelectParticipationsCount()
}

// ParticipantsCount - returns slice of sales ids mapped to participants count
func (p pendingParticipationsQ) ParticipantsCount() ([]core2.SaleIDParticipantsCount, error) {
	return p.offersQ.SelectParticipantsCount()
}
