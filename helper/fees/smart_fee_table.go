package fees

import (
	"math"

	"github.com/saychao/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

//SmartFeeTable is used to built complete fee overview,
//using different level fees

type (
	FeeWrapper struct {
		core.FeeEntry
		NotExists bool
	}
	FeeGroup struct {
		AssetCode string
		FeeType   int
		Subtype   int64
	}
	SmartFeeTable map[FeeGroup][]FeeWrapper
)

func NewSmartFeeTable(fees []core.FeeEntry) (sft SmartFeeTable) {
	sft = SmartFeeTable{}
	for _, entry := range fees {
		key := FeeGroup{
			AssetCode: entry.Asset,
			FeeType:   entry.FeeType,
			Subtype:   entry.Subtype,
		}

		value := FeeWrapper{
			entry,
			false,
		}

		sft[key] = append(sft[key], value)
	}

	return sft
}

func (sft SmartFeeTable) Update(fees []core.FeeEntry) {
	for _, v := range fees {
		key := FeeGroup{
			AssetCode: v.Asset,
			FeeType:   v.FeeType,
			Subtype:   v.Subtype,
		}
		entry := FeeWrapper{
			FeeEntry: v,
		}

		sft[key] = fillFeeGaps(sft[key], entry)
	}
}

func (sft SmartFeeTable) GetValuesByAsset() (byAsset map[string][]FeeWrapper) {
	byAsset = make(map[string][]FeeWrapper)
	for key := range sft {
		byAsset[key.AssetCode] = append(byAsset[key.AssetCode], sft[key]...)
	}

	return byAsset
}

func (sft SmartFeeTable) AddZeroFees(assets []string) {
	for _, asset := range assets {
		for _, ft := range xdr.FeeTypeAll {
			subtypes := []int64{0}
			if ft == xdr.FeeTypePaymentFee {
				subtypes = []int64{
					int64(xdr.PaymentFeeTypeIncoming),
					int64(xdr.PaymentFeeTypeOutgoing),
				}
			}
			for _, st := range subtypes {
				key := FeeGroup{
					AssetCode: asset,
					FeeType:   int(ft),
					Subtype:   st,
				}

				zeroFee := FeeWrapper{
					FeeEntry: core.FeeEntry{
						Asset:      asset,
						Subtype:    st,
						FeeType:    int(ft),
						LowerBound: 0,
						UpperBound: math.MaxInt64,
					},
					NotExists: true,
				}

				sft[key] = fillFeeGaps(sft[key], zeroFee)
			}
		}
	}
}
