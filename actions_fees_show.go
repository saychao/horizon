package horizon

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/exchange"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

// This file contains the actions:
//
// FeesShowAction: renders fees for operationType
type FeesShowAction struct {
	Action
	converter *exchange.Converter

	FeeType   int
	Asset     string
	Subtype   int64
	AccountID string
	Account   *core.Account

	Amount string

	Fee regources.FeeEntry
}

// JSON is a method for actions.JSON
func (action *FeesShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.createConverter,
		action.loadData,
		func() {
			hal.Render(action.W, action.Fee)
		},
	)
}
func (action *FeesShowAction) loadParams() {
	action.FeeType = int(action.GetInt32("fee_type"))
	action.Asset = action.GetNonEmptyString("asset")
	action.Subtype = action.GetInt64("subtype")
	action.AccountID = action.GetString("account")
	action.Amount = action.GetString("amount")
}

func (action *FeesShowAction) createConverter() {
	var err error
	action.converter, err = action.CreateConverter()
	if err != nil {
		action.Log.WithError(err).Error("Failed to init converter")
		action.Err = &problem.ServerError
		return
	}
}

func (action *FeesShowAction) loadData() {
	var err error
	if action.AccountID != "" {
		action.Account, err = action.CoreQ().Accounts().ByAddress(action.AccountID)
	}

	if err != nil {
		action.Log.WithError(err).Error("Failed to load account to get fee")
		action.Err = &problem.ServerError
		return
	}

	am := int64(0)
	if action.Amount != "" {
		amXdr, _ := amount.Parse(action.Amount)
		am = int64(amXdr)
	}
	result, err := action.CoreQ().FeeByTypeAssetAccount(action.FeeType, action.Asset, action.Subtype, action.Account, am)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load fee by asset and type")
		action.Err = &problem.ServerError
		return
	}

	if result == nil {
		result = new(core.FeeEntry)
		result.Asset = action.Asset
		result.FeeType = action.FeeType
	}

	percentFee := action.GetPercentFee(result.Percent)
	if action.Err != nil {
		return
	}

	result.Percent = percentFee

	action.Fee = resource.NewFeeEntry(*result)
}

func (action *FeesShowAction) GetPercentFee(percentFee int64) int64 {
	// request does not require to calculate
	if action.Amount == "" {
		return percentFee
	}

	am, err := amount.Parse(action.Amount)
	if err != nil {
		action.SetInvalidField("amount", errors.Wrap(err, "failed to parse"))
		return 0
	}

	asset, err := action.CoreQ().Assets().ByCode(action.Asset)
	if err != nil {
		action.Log.WithError(err).Error("failed to load asset from core db")
		action.Err = &problem.ServerError
		return 0
	}

	if asset == nil {
		action.SetInvalidField("asset", errors.From(errors.New("does not exist"), logan.F{
			"asset": action.Asset,
		}))
		return 0
	}

	res, isOverflow := action.CalculatePercentFee(percentFee, am, asset.GetMinimumAmount())
	if isOverflow {
		action.SetInvalidField("amount", errors.New("is too big - overflow"))
		return 0
	}

	return res
}
