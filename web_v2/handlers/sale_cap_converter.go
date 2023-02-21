package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/exchange"
	"github.com/saychao/horizon/web_v2/ctx"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	regources "gitlab.com/tokend/regources/generated"
)

// saleCapConverter - helper struct to populate current caps for the sale
type saleCapConverter struct {
	converter *exchange.Converter
}

// newSaleCapConverterForHandler - creates new instance of saleCapConverter.
// returns nil and renders server error if failed to create
func newSaleCapConverterForHandler(w http.ResponseWriter, r *http.Request) *saleCapConverter {
	repo := ctx.CoreRepo(r)
	assetsProvider := struct {
		core2.AssetsQ
		core2.AssetPairsQ
	}{
		AssetsQ:     core2.NewAssetsQ(repo),
		AssetPairsQ: core2.NewAssetPairsQ(repo),
	}

	converter, err := exchange.NewConverter(assetsProvider)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to create new converter")
		ape.Render(w, problems.InternalError())
		return nil
	}

	return newSaleCapConverter(converter)
}

// newSaleCapConverter - creates new instance of saleCapConverter
func newSaleCapConverter(converter *exchange.Converter) *saleCapConverter {
	return &saleCapConverter{
		converter: converter,
	}
}

// PopulateSaleCap - populates current caps of specified sale.
// returns error if failed to populate cap or it's not possible to calculate it
func (c *saleCapConverter) PopulateSaleCap(sale *history2.Sale) error {

	totalCapInDefaultQuoteAsset, err := c.getCurrentCapInDefaultQuote(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to get current cap in default quote asset", logan.F{
			"sale_id": sale.ID,
		})
	}

	sale.CurrentCap = regources.Amount(totalCapInDefaultQuoteAsset)

	for i := range sale.QuoteAssets.QuoteAssets {
		quoteAsset := &sale.QuoteAssets.QuoteAssets[i]
		err = c.populateQuoteAssetCap(quoteAsset, sale, totalCapInDefaultQuoteAsset)
		if err != nil {
			return errors.Wrap(err, "failed to populate quote asset cap", logan.F{
				"quote_asset": quoteAsset.Asset,
				"sale_id":     sale.ID,
			})
		}
	}

	return nil
}

func (c *saleCapConverter) populateQuoteAssetCap(quoteAsset *history2.SaleQuoteAsset, sale *history2.Sale,
	totalCapInDefaultQuoteAsset int64) error {

	totalCapInQuote, err := c.converter.TryToConvertWithOneHop(totalCapInDefaultQuoteAsset, sale.DefaultQuoteAsset,
		quoteAsset.Asset)
	if err != nil {
		return errors.Wrap(err, "failed to convert total cap in default to quote")
	}

	if totalCapInQuote == nil {
		return errors.New("failed to convert total cap in default to quote: failed to find path")
	}

	quoteAsset.TotalCurrentCap = regources.Amount(*totalCapInQuote)

	var hardCapInQuote *int64
	hardCapInQuote, err = c.converter.TryToConvertWithOneHop(int64(sale.HardCap), sale.DefaultQuoteAsset, quoteAsset.Asset)
	if err != nil {
		return errors.Wrap(err, "failed to convert hard cap")
	}

	if hardCapInQuote == nil {
		return errors.New("failed to convert hard cap to quote asset")
	}

	quoteAsset.HardCap = regources.Amount(*hardCapInQuote)
	return nil
}

// PopulateSalesCaps - populates current caps of specified sales
// returns error if failed to populate cap or it's not possible to calculate it
func (c *saleCapConverter) PopulateSalesCaps(sales []history2.Sale) error {
	for i := range sales {
		err := c.PopulateSaleCap(&sales[i])
		if err != nil {
			return errors.Wrap(err, "failed to populate cap for sale", logan.F{
				"sale_id": sales[i].ID,
			})
		}
	}

	return nil
}

func (c *saleCapConverter) getCurrentCapInDefaultQuote(sale history2.Sale) (int64, error) {
	var totalCapInDefaultQuoteAsset int64
	for _, quoteAsset := range sale.QuoteAssets.QuoteAssets {

		currentCapInDefaultQuoteAsset, err := c.converter.TryToConvertWithOneHop(int64(quoteAsset.CurrentCap),
			quoteAsset.Asset, sale.DefaultQuoteAsset)
		if err != nil {
			return 0, errors.Wrap(err, "failed to convert current cap to default quote asset", logan.F{
				"asset":               quoteAsset.Asset,
				"default_quote_asset": sale.DefaultQuoteAsset,
			})
		}

		if currentCapInDefaultQuoteAsset == nil {
			return 0, errors.From(errors.New("failed to convert to current cap: no path found or overflow"), logan.F{
				"asset":               quoteAsset.Asset,
				"default_quote_asset": sale.DefaultQuoteAsset,
			})
		}

		var isOk bool
		totalCapInDefaultQuoteAsset, isOk = amount.SafePositiveSum(totalCapInDefaultQuoteAsset, *currentCapInDefaultQuoteAsset)
		if !isOk {
			return 0, errors.From(errors.New("failed to find total cap in default quote asset: "+
				"overflow of total cap"), logan.F{
				"sale_id": sale.ID,
			})
		}
	}

	return totalCapInDefaultQuoteAsset, nil
}
