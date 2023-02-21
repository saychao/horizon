package horizon

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/exchange"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

type AccountDetailedBalancesAction struct {
	Action
	converter *exchange.Converter

	AccountID string

	Balances []core.Balance
	Assets   []core.Asset
	Sales    []history.Sale

	AssetCodes []string

	Resource []resource.Balance
}

func (action *AccountDetailedBalancesAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadBalances,
		action.groupBalancesByAsset,
		action.loadAssets,
		action.createConverter,
		action.loadSales,
		action.loadResource,
		func() {
			if action.Resource == nil {
				action.Resource = make([]resource.Balance, 0)
			}

			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountDetailedBalancesAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *AccountDetailedBalancesAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *AccountDetailedBalancesAction) getConvertToAsset(fieldName string) string {
	code := action.GetString(fieldName)
	if code != "" {
		return code
	}

	statsQuoteAsset, err := action.CoreQ().Assets().ForPolicy(uint32(xdr.AssetPolicyStatsQuoteAsset)).Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load stats quote asset")
		action.Err = &problem.ServerError
		return ""
	}

	if len(statsQuoteAsset) == 0 {
		action.SetInvalidField(fieldName, errors.New("stats quote asset is not specified. Explicitly specify convert_to asset or create stats quote asset"))
		return ""
	}

	if len(statsQuoteAsset) > 1 {
		action.Log.Error("unexpected number of stats quote assets. Expected [0,1]")
		action.Err = &problem.ServerError
		return ""
	}

	return statsQuoteAsset[0].Code

}

func (action *AccountDetailedBalancesAction) loadBalances() {
	var err error
	action.Balances, err = action.CoreQ().Balances().ByAddress(action.AccountID).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load balances")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) groupBalancesByAsset() {
	assetsMap := map[string]bool{}
	for _, balance := range action.Balances {
		if _, ok := assetsMap[balance.Asset]; ok {
			continue
		}

		assetsMap[balance.Asset] = true
		action.AssetCodes = append(action.AssetCodes, balance.Asset)
	}
}

func (action *AccountDetailedBalancesAction) loadAssets() {
	if len(action.AssetCodes) == 0 {
		return
	}

	var err error
	action.Assets, err = action.CoreQ().Assets().ForCodes(action.AssetCodes).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load assets for codes")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) createConverter() {
	var err error
	action.converter, err = action.CreateConverter()
	if err != nil {
		action.Log.WithError(err).Error("Failed to init converter")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) loadSales() {
	if len(action.AssetCodes) == 0 {
		return
	}

	var err error
	action.Sales, err = selectSalesWithCurrentCap(action.HistoryQ().Sales().ForBaseAssets(action.AssetCodes...), action.converter)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load sales")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountDetailedBalancesAction) loadResource() {
	var convertToAsset string
	if len(action.Balances) != 0 {
		convertToAsset = action.getConvertToAsset("convert_to")
		if action.Err != nil {
			return
		}
	}

	for _, record := range action.Balances {
		var r resource.Balance
		r.Populate(record)
		asset := findAssetByAssetCode(record.Asset, action.Assets)
		if asset == nil {
			action.Log.WithField("asset_code", record.Asset).Error("Failed to find asset for existing balance")
			action.Err = &problem.ServerError
			return
		}

		r.AssetDetails = asset
		var err error
		r.AssetDetails.Sales, err = findAllSalesForAsset(action.CoreQ(), asset.Code, action.Sales)
		if err != nil {
			action.Log.WithError(err).Error("failed to find all sales for asset")
			action.Err = &problem.ServerError
			return
		}

		r.ConvertedBalance, err = convertAmount(record.Amount, r.Asset, convertToAsset, action.converter)
		if err != nil {
			action.Log.WithError(err).Error("Failed to convert balance")
			action.Err = &problem.ServerError
			return
		}

		r.ConvertedLocked, err = convertAmount(record.Locked, r.Asset, convertToAsset, action.converter)
		if err != nil {
			action.Log.WithError(err).Error("failed to convert locked amount")
			action.Err = &problem.ServerError
			return
		}

		r.ConvertedToAsset = convertToAsset
		action.Resource = append(action.Resource, r)
	}
}

func convertAmount(balance int64, fromAsset, toAsset string, converter *exchange.Converter) (string, error) {
	convertedAmount, err := converter.TryToConvertWithOneHop(balance, fromAsset, toAsset)
	if err != nil {
		return "", err
	}

	if convertedAmount == nil {
		return amount.String(0), nil
	}

	return amount.String(*convertedAmount), nil
}

func findAssetByAssetCode(code string, assets []core.Asset) *resource.Asset {
	for i := range assets {
		if code != assets[i].Code {
			continue
		}

		var result resource.Asset
		result.Populate(&assets[i])
		return &result
	}

	return nil
}

func findAllSalesForAsset(q core.QInterface, code string, sales []history.Sale) ([]resource.Sale, error) {
	var result []resource.Sale
	for i := range sales {
		if sales[i].BaseAsset != code {
			continue
		}

		var sale resource.Sale
		sale.Populate(&sales[i])

		err := populateSaleWithStats(sales[i].ID, &sale, q)
		if err != nil {
			return nil, errors.Wrap(err, "failed to populate sale with stats")
		}
		result = append(result, sale)
	}

	return result, nil
}
