package ingest

import (
	"strconv"
	"time"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/ingest/ingestion"
	"github.com/saychao/horizon/log"
	"gitlab.com/tokend/go/amount"
)

type assetPairPriceKey struct {
	BaseAsset  string
	QuoteAsset string
}

func newAssetPairPriceKey(baseAsset, quoteAsset string) assetPairPriceKey {
	return assetPairPriceKey{
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
	}
}

type PriceHistoryProvider struct {
	prices          map[assetPairPriceKey]int64
	ledgerCloseTime time.Time

	log *log.Entry
}

func (h *PriceHistoryProvider) Init(assetPairs []core.AssetPair, ledgerCloseTime time.Time) {
	h.ledgerCloseTime = ledgerCloseTime
	h.log = log.WithField("service", "price_history_provider")
	h.prices = make(map[assetPairPriceKey]int64)

	for _, assetPair := range assetPairs {
		if assetPair.BaseAsset == assetPair.QuoteAsset {
			continue
		}

		h.Put(assetPair.BaseAsset, assetPair.QuoteAsset, assetPair.CurrentPrice)
	}
}

func (h *PriceHistoryProvider) Put(base, quote string, price int64) {
	h.prices[newAssetPairPriceKey(base, quote)] = price
}

func (h *PriceHistoryProvider) ToPricePoints() ([]ingestion.LedgerPricePoint, error) {
	result := make([]ingestion.LedgerPricePoint, 0, len(h.prices))

	for assetPair, pricePoint := range h.prices {
		price, err := strconv.ParseFloat(amount.String(pricePoint), 64)
		if err != nil {
			h.log.WithError(err).Error("Failed to get price history")
			return nil, err
		}

		result = append(result, ingestion.LedgerPricePoint{
			BaseAsset:  assetPair.BaseAsset,
			QuoteAsset: assetPair.QuoteAsset,
			PricePoint: history.PricePoint{
				Price:     price,
				Timestamp: h.ledgerCloseTime,
			},
		})
	}

	return result, nil
}
