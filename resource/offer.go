package resource

import (
	"strconv"

	"github.com/saychao/horizon/db2/history"

	"time"

	"github.com/saychao/horizon/db2/core"
	"gitlab.com/tokend/regources"
)

func PopulateOffer(o core.Offer) regources.Offer {
	return regources.Offer{
		OwnerID:        o.OwnerID,
		OfferID:        o.OfferID,
		OrderBookID:    o.OrderBookID,
		OfferData:      PopulateOfferData(o),
		BaseBalanceID:  o.BaseBalanceID,
		QuoteBalanceID: o.QuoteBalanceID,
		Fee:            regources.Amount(o.Fee),
		PT:             strconv.FormatUint(o.OfferID, 10),
	}
}

func PopulateOfferData(o core.Offer) regources.OfferData {
	return regources.OfferData{
		BaseAssetCode:  o.BaseAssetCode,
		QuoteAssetCode: o.QuoteAssetCode,
		IsBuy:          o.IsBuy,
		BaseAmount:     regources.Amount(o.BaseAmount),
		QuoteAmount:    regources.Amount(o.QuoteAmount),
		Price:          regources.Amount(o.Price),
		CreatedAt:      regources.Time(time.Unix(o.CreatedAt, 0).UTC()),
	}
}

func PopulateHistoryOffer(histOffer history.Offer) regources.HistoryOffer {
	offerID := ""
	if histOffer.OfferID != 0 {
		offerID = strconv.FormatInt(histOffer.OfferID, 10)
	}

	return regources.HistoryOffer{
		PT:                strconv.FormatInt(histOffer.ID, 10),
		OfferID:           offerID,
		OwnerID:           histOffer.OwnerID,
		BaseAsset:         histOffer.BaseAsset,
		QuoteAsset:        histOffer.QuoteAsset,
		IsBuy:             histOffer.IsBuy,
		InitialBaseAmount: regources.Amount(histOffer.InitialBaseAmount),
		CurrentBaseAmount: regources.Amount(histOffer.CurrentBaseAmount),
		Price:             regources.Amount(histOffer.Price),
		IsCanceled:        histOffer.IsCanceled,
		CreatedAt:         regources.Time(histOffer.CreatedAt),
	}
}
