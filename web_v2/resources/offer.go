package resources

import (
	"fmt"
	"time"

	"github.com/SafeRE-IT/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewOffer creates new instance of Offer from provided one
func NewOffer(record core2.Offer) regources.Offer {
	return regources.Offer{
		Key: regources.Key{
			ID:   fmt.Sprint(record.OfferID),
			Type: regources.OFFERS,
		},
		Attributes: regources.OfferAttributes{
			IsBuy:       record.IsBuy,
			OrderBookId: fmt.Sprint(record.OrderBookID),
			CreatedAt:   time.Unix(record.CreatedAt, 0).UTC(),
			BaseAmount:  regources.Amount(record.BaseAmount),
			QuoteAmount: regources.Amount(record.QuoteAmount),
			Price:       regources.Amount(record.Price),
			Fee: regources.Fee{
				CalculatedPercent: regources.Amount(record.Fee),
			},
		},
	}
}
