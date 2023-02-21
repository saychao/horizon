package horizon

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/exchange"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
)

// SaleShowAction renders a sale found by its ID.
type SaleShowAction struct {
	Action
	RequestID uint64
	Record    *history.Sale
	result    resource.Sale
}

// JSON is a method for actions.JSON
func (action *SaleShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		action.populateResult,
		func() {
			hal.Render(action.W, action.result)
		},
	)
}

func (action *SaleShowAction) loadParams() {
	action.RequestID = action.GetUInt64("id")
}

func (action *SaleShowAction) loadRecord() {
	var err error
	action.Record, err = action.HistoryQ().Sales().ByID(action.RequestID)
	if err != nil {
		action.Log.WithError(err).
			WithField("request_id", action.RequestID).
			Error("failed to load sale")
		action.Err = &problem.ServerError
		return
	}

	if action.Record == nil {
		action.Err = &problem.NotFound
		return
	}

	err = action.populateTotalCurrentCap()
	if err != nil {
		action.Log.WithError(err).Error("failed to populate total current cap")
		action.Err = &problem.ServerError
		return
	}
}

func (action *SaleShowAction) populateTotalCurrentCap() error {
	converter, err := action.CreateConverter()
	if err != nil {
		return errors.Wrap(err, "failed to init converter")
	}

	totalCapInDefaultQuoteAsset, err := getCurrentCapInDefaultQuote(*action.Record, converter)
	if err != nil {
		return errors.Wrap(err, "failed to get current cap in default quote asset")
	}

	action.Record.CurrentCap = amount.String(totalCapInDefaultQuoteAsset)

	for i := range action.Record.QuoteAssets.QuoteAssets {
		quoteAsset := &action.Record.QuoteAssets.QuoteAssets[i]
		totalCapInQuote, err := converter.TryToConvertWithOneHop(totalCapInDefaultQuoteAsset, action.Record.DefaultQuoteAsset, quoteAsset.Asset)
		if err != nil {
			return errors.Wrap(err, "failed to convert total cap in default to quote")
		}

		if totalCapInQuote == nil {
			return errors.New("failed to convert total cap in default to quote: failed to find path")
		}

		quoteAsset.TotalCurrentCap = amount.String(*totalCapInQuote)

		hardCapInQuote, err := converter.TryToConvertWithOneHop(int64(action.Record.HardCap), action.Record.DefaultQuoteAsset, quoteAsset.Asset)
		if err != nil {
			return errors.Wrap(err, "failed to convert hard cap")
		}

		if hardCapInQuote == nil {
			return errors.New("failed to convert hard cap to quote asset")
		}

		quoteAsset.HardCap = amount.String(*hardCapInQuote)
	}

	return nil
}

func (action *SaleShowAction) populateResult() {
	action.result.Populate(action.Record)
	err := populateSaleWithStats(action.Record.ID, &action.result, action.CoreQ())
	if err != nil {
		action.Log.WithError(err).
			WithField("request_id", action.RequestID).
			Error("failed to populate stat for sale")
		action.Err = &problem.ServerError
		return
	}
}

func populateSaleWithStats(saleID uint64, sale *resource.Sale, q core.QInterface) error {
	var offers []core.Offer
	err := q.Offers().ForOrderBookID(saleID).Select(&offers)
	if err != nil {
		return errors.Wrap(err, "failed to load offers for sale", map[string]interface{}{
			"sale_id": saleID,
		})
	}

	balances, err := q.Balances().ByAsset(sale.BaseAsset).Select()
	if err != nil {
		return errors.Wrap(err, "failed to load balances for sale's base asset", map[string]interface{}{
			"sale_id": saleID,
		})
	}

	err = sale.PopulateStat(offers, balances)
	if err != nil {
		return errors.Wrap(err, "failed to populate stats for sale", map[string]interface{}{
			"sale_id": saleID,
		})
	}

	return nil
}

func selectSalesWithCurrentCap(q history.SalesQ, converter *exchange.Converter) ([]history.Sale, error) {
	result, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sales")
	}

	for i := range result {
		currentCapInDefaultQuote, err := getCurrentCapInDefaultQuote(result[i], converter)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate current cap in default quote")
		}

		result[i].CurrentCap = amount.String(currentCapInDefaultQuote)
	}

	return result, nil
}

func getCurrentCapInDefaultQuote(sale history.Sale, converter *exchange.Converter) (int64, error) {
	totalCapInDefaultQuoteAsset := int64(0)
	for _, quoteAsset := range sale.QuoteAssets.QuoteAssets {
		currentCap, err := amount.Parse(quoteAsset.CurrentCap)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse current cap")
		}

		currentCapInDefaultQuoteAsset, err := converter.TryToConvertWithOneHop(currentCap, quoteAsset.Asset, sale.DefaultQuoteAsset)
		if err != nil {
			return 0, errors.Wrap(err, "failed to convert current cap to default quote asset")
		}

		if currentCapInDefaultQuoteAsset == nil {
			return 0, errors.New("failed to convert to current cap: no path found")
		}

		var isOk bool
		totalCapInDefaultQuoteAsset, isOk = amount.SafePositiveSum(totalCapInDefaultQuoteAsset, *currentCapInDefaultQuoteAsset)
		if !isOk {
			return 0, errors.New("failed to find total cap in default quote asset: overflow")
		}
	}

	return totalCapInDefaultQuoteAsset, nil
}
