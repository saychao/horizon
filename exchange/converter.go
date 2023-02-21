package exchange

import (
	core "github.com/saychao/horizon/db2/core2"
	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

// Converter - helper struct which allows to convert amount in asset into amount in another asset
// supports indirect conversion with one hop
type Converter struct {
	assetProvider assetProvider
	baseAssets    []string
}

// NewConverter - creates new ready to work instance of Converter
func NewConverter(assetProvider assetProvider) (*Converter, error) {
	baseAssets, err := assetProvider.SelectByPolicy(uint64(xdr.AssetPolicyBaseAsset))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load base assets")
	}

	result := &Converter{
		assetProvider: assetProvider,
	}

	for i := range baseAssets {
		result.baseAssets = append(result.baseAssets, baseAssets[i].Code)
	}

	return result, nil
}

// TryToConvertWithOneHop - tries to convert amount in fromAsset into amount in toAsset
// if no direct conversion possible, tries to do it with one cop;
// if several conversions are available selects one that results in maximizing converted amount
// if fails to find path to convert - returns nil, nil
func (c *Converter) TryToConvertWithOneHop(amount int64, fromAsset, toAsset string) (*int64, error) {
	// converting to self
	if fromAsset == toAsset {
		return &amount, nil
	}

	directPair, err := c.tryLoadDirect(fromAsset, toAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load direct asset pair")
	}

	if directPair != nil {
		var result int64
		var isConverted bool
		result, isConverted, err = c.convertToDestAsset(*directPair, toAsset, amount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert using direct pair")
		}

		if !isConverted {
			return nil, nil
		}

		return &result, nil
	}

	fromAssetPairs, err := c.loadPairsWithBaseAssets(fromAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load pairs with base asset for from asset")
	}

	toAssetPairs, err := c.loadPairsWithBaseAssets(toAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load pairs with base asset for to asset")
	}

	return c.convertWithMaxPath(amount, fromAsset, toAsset, fromAssetPairs, toAssetPairs)
}

func (c *Converter) loadPairsWithBaseAssets(asset string) ([]core.AssetPair, error) {

	direct, err := c.assetProvider.SelectByAssets([]string{asset}, c.baseAssets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load direct asset pairs")
	}

	reverse, err := c.assetProvider.SelectByAssets(c.baseAssets, []string{asset})
	if err != nil {
		return nil, errors.Wrap(err, "failed to load reverse asset pairs")
	}

	return append(direct, reverse...), nil
}

func (c *Converter) tryLoadDirect(fromAsset, toAsset string) (*core.AssetPair, error) {
	assets := []string{fromAsset, toAsset}
	pairs, err := c.assetProvider.SelectByAssets(assets, assets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load asset pairs by assets")
	}

	pairToFind := core.AssetPair{
		Base:  fromAsset,
		Quote: toAsset,
	}
	for i, pair := range pairs {
		if equalsOrInverted(pair, pairToFind) {
			return &pairs[i], nil
		}
	}

	return nil, nil
}

func (c *Converter) convertWithMaxPath(amount int64, fromAsset, toAsset string, fromPairs, toPairs []core.AssetPair) (*int64, error) {
	converted := false
	var result int64
	for _, fromPair := range fromPairs {
		hopAmount, isConverted, err := c.convertFromSourceAsset(fromPair, fromAsset, amount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert from asset to hop asset")
		}

		if !isConverted {
			continue
		}

		for _, toPair := range toPairs {
			if equalsOrInverted(fromPair, toPair) {
				return nil, errors.From(errors.New("unexpected state: received same pairs"), map[string]interface{}{
					"fromPair": fromPair,
					"toPair":   toPair,
				})
			}
			if !isOverlaps(fromPair, toPair) {
				continue
			}

			var destAmount int64
			destAmount, isConverted, err = c.convertToDestAsset(toPair, toAsset, hopAmount)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert to toAsset")
			}

			if !isConverted {
				continue
			}

			converted = true
			if destAmount > result {
				result = destAmount
			}
		}
	}

	if converted {
		return &result, nil
	}

	return nil, nil
}

// convertToDestAsset - converts specified amount to dest asset using current price,
// returns false - if failed
func (c *Converter) convertToDestAsset(pair core.AssetPair, destCode string, amountToConvert int64) (int64, bool, error) {
	if pair.CurrentPrice == 0 {
		return 0, false, errors.New("Price is invalid")
	}

	destAsset, err := c.assetProvider.GetByCode(destCode)
	if err != nil {
		return 0, false, errors.From(errors.New("failed to select dest asset"), logan.F{
			"destCode": destCode,
		})
	}

	if destAsset == nil {
		return 0, false, errors.From(errors.New("asset not found"), logan.F{
			"destCode": destCode,
		})
	}

	switch destCode {
	case pair.Quote:
		result, isOverflow := amount.BigDivide(amountToConvert, pair.CurrentPrice,
			amount.One, amount.ROUND_UP, destAsset.GetMinimumAmount())
		return result, !isOverflow, nil
	case pair.Base:
		result, isOverflow := amount.BigDivide(amountToConvert, amount.One,
			pair.CurrentPrice, amount.ROUND_UP, destAsset.GetMinimumAmount())
		return result, !isOverflow, nil
	default:
		return 0, false, errors.From(errors.New("unexpected dest code"), logan.F{
			"base":        pair.Base,
			"quote":       pair.Quote,
			"actual dest": destCode,
		})
	}
}

// isOverlaps - returns true if one of the assets of anotherPair is equal to one of the assets of pair
func isOverlaps(pair, anotherPair core.AssetPair) bool {
	return contains(pair, anotherPair.Base) || contains(pair, anotherPair.Quote)
}

// Contains - returns true if base or quote equal to asset
func contains(pair core.AssetPair, asset string) bool {
	return pair.Base == asset || pair.Quote == asset
}

// equalsOrInverted - returns true if other pair is the same or inverted
func equalsOrInverted(pair, other core.AssetPair) bool {
	return (pair.Base == other.Base && pair.Quote == other.Quote) ||
		(pair.Base == other.Quote && pair.Quote == other.Base)
}

// ConvertFromSourceAsset - converts specified amount from source to another asset in pair using current price,
// returns false - if failed
func (c *Converter) convertFromSourceAsset(pair core.AssetPair, sourceCode string, amountToConvert int64) (int64, bool, error) {
	destCode := ""
	switch sourceCode {
	case pair.Base:
		destCode = pair.Quote
	case pair.Quote:
		destCode = pair.Base
	default:
		return 0, false, errors.From(errors.New("unexpected source code"), logan.F{
			"base":          pair.Base,
			"quote":         pair.Quote,
			"actual source": sourceCode,
		})
	}

	return c.convertToDestAsset(pair, destCode, amountToConvert)
}
