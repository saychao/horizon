package handlers

import (
	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	regources "gitlab.com/tokend/regources/generated"
)

type closedParticipationsQ struct {
	participationQ history2.SaleParticipationQ
}

func newClosedParticipationQ(request *requests.GetSaleParticipations, q history2.SaleParticipationQ, sale *history2.Sale,
) closedParticipationsQ {
	q = q.FilterBySaleParams(sale.ID, sale.BaseAsset, sale.OwnerAddress).Page(request.PageParams)

	return closedParticipationsQ{
		participationQ: q,
	}
}

// FilterByParticipant - filters out participations by participant address
func (q closedParticipationsQ) FilterByParticipant(id string) participationsQ {
	q.participationQ = q.participationQ.FilterByParticipant(id)
	return q
}

// FilterByQuoteAsset - filters out participations by quote asset
func (q closedParticipationsQ) FilterByQuoteAsset(code string) participationsQ {
	q.participationQ = q.participationQ.FilterByQuoteAsset(code)
	return q
}

// Select - select records from db and wraps them to participations
func (q closedParticipationsQ) Select() ([]regources.SaleParticipation, error) {
	matches, err := q.participationQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load matches from db")
	}

	result := make([]regources.SaleParticipation, 0, len(matches))
	for _, m := range matches {
		result = append(result, resources.NewSaleParticipation(
			m.ID,
			m.ParticipantID,
			m.BaseAsset,
			m.QuoteAsset,
			amount.MustParseU(m.QuoteAmount),
		))
	}

	return result, nil
}

// ParticipationsCount - returns slice of sales ids mapped to participants count
func (q closedParticipationsQ) ParticipationsCount() ([]core2.SaleIDParticipantsCount, error) {
	return q.participationQ.SelectParticipationsCount()
}

// ParticipantsCount - returns slice of sales ids mapped to participants count
func (q closedParticipationsQ) ParticipantsCount() ([]core2.SaleIDParticipantsCount, error) {
	return q.participationQ.SelectParticipantsCount()
}
