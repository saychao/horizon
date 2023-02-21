package horizon

import (
	"database/sql"

	"github.com/go-errors/errors"
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// This file contains the actions:
//
// FeesAllAction: renders all registration requests

//FeesAllAction show all fees

type FeesAllAction struct {
	Action
	AccountType *uint64
	Account     string

	IsOverview bool

	Response resource.FeesResponse
}

func (action *FeesAllAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *FeesAllAction) loadParams() {
	action.Account = action.GetString("account_id")
	action.AccountType = action.GetOptionalUint64("account_type")
	if action.Account != "" && action.AccountType != nil {
		action.SetInvalidField("account_type", errors.New("It's not allowed to set both filters"))
	}
}

func (action *FeesAllAction) getAccountType(name string) *int32 {
	rawAccountType := action.GetInt32(name)
	if action.Err != nil {
		return nil
	}

	if rawAccountType == 0 {
		return nil
	}

	action.SetInvalidField(name, errors.New("Invalid"))
	return nil
}

func (action *FeesAllAction) loadData() {
	q := action.CoreQ().FeeEntries()
	// for the overview we need to return all the fee rules we have, so we just ignore filters
	if !action.IsOverview {
		// pass all the filters. Q will resolve them correctly
		q = q.ForAccount(action.Account).ForAccountType(action.AccountType)
	}

	actualFees, err := q.Select()
	if err != nil {
		if err != sql.ErrNoRows {
			action.Err = &problem.ServerError
			action.Log.WithStack(err).WithError(err).Error("Could not get fee from the database")
			return
		}

		err = nil
	}

	// convert to map of regourcess
	action.Response.Fees = map[xdr.AssetCode][]regources.FeeEntry{}
	var fee regources.FeeEntry
	for _, coreFee := range actualFees {
		fee = resource.NewFeeEntry(coreFee)
		ac := xdr.AssetCode(coreFee.Asset)
		action.Response.Fees[ac] = append(action.Response.Fees[ac], fee)
	}

	// for overview we do not need to populate default fees
	if action.IsOverview {
		return
	}

	assets, err := action.CoreQ().Assets().Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load assets")
		action.Err = &problem.ServerError
		return
	}

	for _, asset := range assets {
		ac := xdr.AssetCode(asset.Code)
		action.Response.Fees[ac] = action.addDefaultEntriesForAsset(asset, action.Response.Fees[ac])
	}
}

func feesContainsType(feeType int, entries []regources.FeeEntry) bool {
	for _, entry := range entries {
		if entry.FeeType == feeType {
			return true
		}
	}

	return false
}

func (action *FeesAllAction) addDefaultEntriesForAsset(asset core.Asset, entries []regources.FeeEntry) []regources.FeeEntry {
	for _, feeType := range xdr.FeeTypeAll {
		subtypes := []int64{0}
		switch feeType {
		case xdr.FeeTypePaymentFee:
			subtypes = []int64{int64(xdr.PaymentFeeTypeIncoming), int64(xdr.PaymentFeeTypeOutgoing)}
		case xdr.FeeTypeOperationFee:
			subtypes = []int64{}
			for _, subtype := range xdr.OperationTypeAll {
				subtypes = append(subtypes, int64(subtype))
			}
		}

		for _, subtype := range subtypes {
			entries = append(entries, action.getDefaultFee(asset.Code, int(feeType), subtype))
		}
	}

	return entries
}

func (action *FeesAllAction) getDefaultFee(asset string, feeType int, subType int64) regources.FeeEntry {
	accountType := uint64(0)
	if action.AccountType != nil {
		accountType = *action.AccountType
	}

	return regources.FeeEntry{
		Asset:       asset,
		FeeType:     feeType,
		Subtype:     subType,
		Percent:     "0",
		Fixed:       "0",
		LowerBound:  "0",
		UpperBound:  "0",
		AccountRole: accountType,
		AccountID:   action.Account,
		FeeAsset:    asset,
	}
}
