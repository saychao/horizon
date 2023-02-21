package resource

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/helper/fees"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

// Populate fills out the resource's fields
func NewFeeEntry(cfee core.FeeEntry) (fee regources.FeeEntry) {
	fee.FeeType = cfee.FeeType
	fee.Asset = cfee.Asset
	fee.Fixed = amount.String(cfee.Fixed)
	fee.Percent = amount.String(cfee.Percent)
	fee.Subtype = cfee.Subtype
	fee.AccountID = cfee.AccountID
	fee.AccountRole = uint64(cfee.AccountRole)
	fee.LowerBound = amount.String(cfee.LowerBound)
	fee.UpperBound = amount.String(cfee.UpperBound)
	fee.Exists = true
	return fee
}

func NewFeeEntryFromWrapper(wrapper fees.FeeWrapper) (fee regources.FeeEntry) {
	fee = NewFeeEntry(wrapper.FeeEntry)
	fee.Exists = !wrapper.NotExists
	return fee
}
