package horizon

import (
	"github.com/pkg/errors"
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/helper/fees"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// This file contains the actions:
//
// FeesForAccount: renders all the fees for a specific account
//FeesForAccount show all fees for account

type AccountFeesAction struct {
	Action
	Account *core.Account

	Records  fees.SmartFeeTable
	Assets   []string
	Response resource.FeesResponse
}

func (action *AccountFeesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadAssets,
		action.loadFees,
		action.loadResponse,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *AccountFeesAction) loadParams() {
	action.Account = action.GetCoreAccount("id", action.CoreQ())
	if action.Err != nil {
		return
	}
	if action.Account == nil {
		action.SetInvalidField("id", errors.New("Must not be empty"))
		return
	}
}

func (action *AccountFeesAction) loadAssets() {
	accountBalances, err := action.CoreQ().Balances().ByAddress(action.Account.AccountID).Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load balances")
		action.Err = &problem.ServerError
		return
	}

	var assets []string
	for _, balance := range accountBalances {
		assets = append(assets, balance.Asset)
	}
	action.Assets = assets
}

func (action *AccountFeesAction) loadFees() {
	forAccount, err := action.CoreQ().FeeEntries().ForAccount(action.Account.AccountID).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("can't get account fees from the database")
		return
	}

	forAccountType, err := action.CoreQ().FeeEntries().ForAccountType(&action.Account.RoleID).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("can't get account type fees from the database")
		return

	}

	//get general fees set for all, not to be confused with fees for General Account Type
	generalFees, err := action.CoreQ().FeeEntries().ForAccountType(nil).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("can't get general fees from the database")
		return
	}

	sft := fees.NewSmartFeeTable(forAccount)
	sft.Update(forAccountType)
	sft.Update(generalFees)
	sft.AddZeroFees(action.Assets)
	action.Records = sft
}

func (action *AccountFeesAction) loadResponse() {
	byAssets := action.Records.GetValuesByAsset()
	action.Response.Fees = make(map[xdr.AssetCode][]regources.FeeEntry)
	for _, feesForAsset := range byAssets {
		for _, wrapper := range feesForAsset {
			fee := resource.NewFeeEntryFromWrapper(wrapper)
			ac := xdr.AssetCode(wrapper.Asset)
			action.Response.Fees[ac] = append(action.Response.Fees[ac], fee)
		}
	}
}
