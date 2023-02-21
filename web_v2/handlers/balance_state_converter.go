package handlers

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/exchange"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

// balanceStateConverter - helper struct to convert balance states to different assets
type balanceStateConverter struct {
	converter *exchange.Converter
}

// newBalanceStateConverterForHandler - creates new instance of balanceStateConverter.
// returns nil and renders server error if failed to create
func newBalanceStateConverterForHandler(coreRepo *pgdb.DB) (*balanceStateConverter, error) {
	assetsProvider := struct {
		core2.AssetsQ
		core2.AssetPairsQ
	}{
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		AssetPairsQ: core2.NewAssetPairsQ(coreRepo),
	}

	converter, err := exchange.NewConverter(assetsProvider)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new converter")
	}

	return newBalanceStateConverted(converter), nil
}

// newBalanceStateConverted - creates new instance of balanceStateConverter
func newBalanceStateConverted(converter *exchange.Converter) *balanceStateConverter {
	return &balanceStateConverter{
		converter: converter,
	}
}

// Convert - returns converted balance state from existing one
func (c *balanceStateConverter) Convert(balance core2.Balance, toAsset string) (*regources.ConvertedBalanceState, error) {
	convertedAvailable, err := c.converter.TryToConvertWithOneHop(
		int64(balance.Amount),
		balance.AssetCode,
		toAsset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert available amount")
	}

	convertedLocked, err := c.converter.TryToConvertWithOneHop(
		int64(balance.Locked),
		balance.AssetCode,
		toAsset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert locked amount")
	}

	if convertedAvailable == nil || convertedLocked == nil {
		convertedState := resources.NewConvertedBalanceState(
			balance,
			regources.Amount(0),
			regources.Amount(0),
			regources.Amount(0),
			false,
		)

		return &convertedState, nil
	}

	price, err := getPrice(*convertedAvailable, *convertedLocked, int64(balance.Amount), int64(balance.Locked))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get price")
	}

	convertedState := resources.NewConvertedBalanceState(
		balance,
		regources.Amount(*convertedAvailable),
		regources.Amount(*convertedLocked),
		regources.Amount(price),
		true,
	)

	return &convertedState, nil
}

// minimum amount is always 1  because method GetMinimumAmount() returns 10^N where N is 0
func getPrice(convertedAvailable, convertedLocked, balanceAmount, balanceLocked int64) (int64, error) {
	if (convertedAvailable != 0) && (balanceAmount != 0) {
		return get(convertedAvailable, balanceAmount, 1)
	}

	if (convertedAvailable == 0) &&
		(convertedLocked != 0) &&
		(balanceLocked != 0) {
		return get(convertedLocked, balanceLocked, 1)
	}

	return 0, nil
}

func get(quote, base, minimumAmount int64) (int64, error) {
	price, isOverflow := amount.BigDivide(quote, amount.One, base, amount.ROUND_UP, minimumAmount)
	if isOverflow {
		return 0, errors.New("price is overflow")
	}
	return price, nil
}
