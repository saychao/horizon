package fees

import (
	"testing"

	"math"

	"github.com/saychao/horizon/db2/core"
	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

func TestNewSmartFeeTable(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		expected := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []FeeWrapper{
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 5,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 10,
						UpperBound: 15,
					},
				},
			},
		}
		sft := NewSmartFeeTable(fees)
		assert.Equal(t, sft, expected)
	})
}

func TestSmartFeeTable_GetValuesByAsset(t *testing.T) {
	fees := []core.FeeEntry{
		{
			FeeType:    int(xdr.FeeTypePaymentFee),
			Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
			Asset:      "USD",
			LowerBound: 0,
			UpperBound: 5,
		},
		{
			FeeType:    int(xdr.FeeTypePaymentFee),
			Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
			Asset:      "USD",
			LowerBound: 10,
			UpperBound: 15,
		},
	}
	expected := map[string][]FeeWrapper{
		"USD": {
			{
				FeeEntry: core.FeeEntry{
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					FeeType:    int(xdr.FeeTypePaymentFee),
					Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
					Asset:      "USD",
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		},
	}
	sft := NewSmartFeeTable(fees)
	got := sft.GetValuesByAsset()
	assert.Equal(t, got, expected)

}

func TestSmartFeeTable_Update(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []FeeWrapper{
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 5,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 6,
						UpperBound: 9,
						Percent:    2,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 10,
						UpperBound: 15,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 16,
						UpperBound: 20,
						Percent:    2,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 21,
						UpperBound: 25,
						Percent:    3,
					},
				},
			},
		}

		sft := NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("no primary fees", func(t *testing.T) {
		var primaryFees []core.FeeEntry
		secondaryFees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []FeeWrapper{
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 2,
						UpperBound: 20,
						Percent:    2,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 21,
						UpperBound: 25,
						Percent:    3,
					},
				},
			},
		}

		sft := NewSmartFeeTable(primaryFees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("different attributes of secondary fees", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 5,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 10,
				UpperBound: 15,
			},
		}
		secondaryFees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypeOfferFee),
				Subtype:    int64(0),
				Asset:      "USD",
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
			{
				FeeType:    int(xdr.FeeTypeOfferFee),
				Subtype:    int64(0),
				Asset:      "USD",
				LowerBound: 21,
				UpperBound: 25,
				Percent:    3,
			},
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []FeeWrapper{
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 5,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 10,
						UpperBound: 15,
					},
				},
			},

			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypeOfferFee),
				Subtype:   int64(0),
			}: []FeeWrapper{
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypeOfferFee),
						Subtype:    int64(0),
						Asset:      "USD",
						LowerBound: 2,
						UpperBound: 20,
						Percent:    2,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypeOfferFee),
						Subtype:    int64(0),
						Asset:      "USD",
						LowerBound: 21,
						UpperBound: 25,
						Percent:    3,
					},
				},
			},
		}

		sft := NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		assert.Equal(t, sft, expectedFeeTable)
	})
	t.Run("both nil", func(t *testing.T) {
		var primaryFees []core.FeeEntry
		var secondaryFees []core.FeeEntry
		sft := NewSmartFeeTable(primaryFees)
		sft.Update(secondaryFees)
		assert.Empty(t, sft)
	})

	t.Run("valid", func(t *testing.T) {
		fees := []core.FeeEntry{
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 0,
				UpperBound: 815 * amount.One,
			},
			{
				FeeType:    int(xdr.FeeTypePaymentFee),
				Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
				Asset:      "USD",
				LowerBound: 816 * amount.One,
				UpperBound: 1000 * amount.One,
			},
		}
		secondaryFees := []core.FeeEntry{}
		expectedFeeGroup := FeeGroup{
			AssetCode: "USD",
			FeeType:   int(xdr.FeeTypePaymentFee),
			Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
		}
		expectedFeeTable := SmartFeeTable{
			FeeGroup{
				AssetCode: "USD",
				FeeType:   int(xdr.FeeTypePaymentFee),
				Subtype:   int64(xdr.PaymentFeeTypeOutgoing),
			}: []FeeWrapper{
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 0,
						UpperBound: 815 * amount.One,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 815*amount.One + 1,
						UpperBound: 816*amount.One - 1,
					},
					NotExists: true,
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 816 * amount.One,
						UpperBound: 1000 * amount.One,
					},
				},
				{
					FeeEntry: core.FeeEntry{
						FeeType:    int(xdr.FeeTypePaymentFee),
						Subtype:    int64(xdr.PaymentFeeTypeOutgoing),
						Asset:      "USD",
						LowerBound: 1000*amount.One + 1,
						UpperBound: math.MaxInt64,
					},
					NotExists: true,
				},
			},
		}

		sft := NewSmartFeeTable(fees)
		sft.Update(secondaryFees)
		sft.AddZeroFees([]string{"USD"})

		assert.Equal(t, sft[expectedFeeGroup], expectedFeeTable[expectedFeeGroup])
	})
}
