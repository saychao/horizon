package horizon

import (
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"gitlab.com/tokend/regources"
)

type BalancesReportAction struct {
	Action
	threshold int64
	assetType string
	Records   regources.BalancesReport
}

func (action *BalancesReportAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		func() {
			hal.Render(action.W, &action.Records)
		},
	)
}

func (action *BalancesReportAction) loadParams() {
	action.threshold = action.GetPositiveInt64("threshold")
	action.assetType = action.GetNonEmptyString("asset_code")
}

func (action *BalancesReportAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *BalancesReportAction) loadRecords() {
	// find positive balances
	balances, err := action.CoreQ().Balances().NonZero().Select()
	if err != nil {
		action.Log.WithError(err).Warn("failed to filter non-zero balances")
		action.Err = &problem.ServerError
		return
	}

	AmountFromAccount := make(map[string]int64)

	converter, err := action.CreateConverter()
	if err != nil {
		action.Log.Warn("failed to create converter")
		action.Err = &problem.ServerError
		return
	}

	// group balances by account and assign converted amounts
	for _, balance := range balances {
		amount, err := converter.TryToConvertWithOneHop(balance.Amount+balance.Locked, balance.Asset, action.assetType)
		if err != nil {
			action.Log.Warn("failed to convert to asset type, skipping")
			action.Err = &problem.ServerError
			return
		}
		if amount == nil {
			action.Log.WithField("base", balance.Asset).WithField("quote", action.assetType).Warn("the asset is unconvertable, skipping")
			action.Err = &problem.ServerError
			return
		} else {
			AmountFromAccount[balance.AccountID] += *amount
		}
	}

	// filter and count by threshold
	for _, amount := range AmountFromAccount {
		if amount < action.threshold {
			action.Records.TotalAccountsCount.BelowThreshold += 1
		} else {
			action.Records.TotalAccountsCount.AboveThreshold += 1
		}
	}

	// find zero balances
	zeroBalances, err := action.CoreQ().Balances().Zero().Select()
	if err != nil {
		action.Log.WithError(err).Warn("failed to filter zero balances")
		action.Err = &problem.ServerError
		return
	}

	isZeroAmountAccount := make(map[string]bool)

	// group zero balances by account
	for _, balance := range zeroBalances {
		isZeroAmountAccount[balance.AccountID] = true
	}

	action.Records.TotalAccountsCount.ZeroBalance = len(isZeroAmountAccount)
	action.Records.TotalAccountsCount.PositiveBalance = len(AmountFromAccount)
}
