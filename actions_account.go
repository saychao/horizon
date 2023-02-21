package horizon

import (
	"time"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/tokend/regources"
)

// This file contains the actions:
//
// AccountShowAction: details for single account (including stellar-core state)

// AccountShowAction renders a account summary found by its address.
type AccountShowAction struct {
	Action
	Address  string
	Resource resource.Account
}

// JSON is a method for actions.JSON
func (action *AccountShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadBalances,
		action.loadExternalSystemAccountIDs,
		action.loadReferrals,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountShowAction) loadParams() {
	action.Address = action.GetString("id")
}

func (action *AccountShowAction) checkAllowed() {
	action.IsAllowed(action.Address)
}

func (action *AccountShowAction) loadRecord() {
	coreRecord, err := action.CoreQ().
		Accounts().
		ForAddresses(action.Address).
		WithAccountKYC().
		First()

	if err != nil {
		action.Log.WithError(err).Error("Failed to get account from core DB")
		action.Err = &problem.ServerError
		return
	}

	if coreRecord == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Resource.Populate(action.Ctx, *coreRecord)
}

func (action *AccountShowAction) loadBalances() {
	var balances []core.Balance
	err := action.CoreQ().
		BalancesByAddress(&balances, action.Address)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get balances for account")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.SetBalances(balances)
}

func (action *AccountShowAction) loadExternalSystemAccountIDs() {
	exSysIDs, err := action.CoreQ().ExternalSystemAccountID().ForAccount(action.Address).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load external system account ids")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.ExternalSystemAccounts = make([]regources.ExternalSystemAccountID, 0, len(exSysIDs))
	for i := range exSysIDs {
		if exSysIDs[i].ExpiresAt == nil || *exSysIDs[i].ExpiresAt >= time.Now().Unix() {
			result := resource.PopulateExternalSystemAccountID(exSysIDs[i])
			action.Resource.ExternalSystemAccounts = append(action.Resource.ExternalSystemAccounts, result)
		}
	}
}

func (action *AccountShowAction) loadReferrals() {
	var coreReferrals []core.Account
	err := action.CoreQ().Accounts().ForReferrer(action.Address).Select(&coreReferrals)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load referrals")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.Referrals = make([]resource.Referral, len(coreReferrals))
	for i := range coreReferrals {
		action.Resource.Referrals[i].Populate(coreReferrals[i])
	}
}
