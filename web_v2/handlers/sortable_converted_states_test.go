package handlers

import (
	"github.com/magiconair/properties/assert"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources/generated"
	"testing"
)

type sortEntry struct {
	convertedAvailable string
	initialAvailable   string
	isConverted        bool
}

func fromSortEntries(states []sortEntry) []regources.ConvertedBalanceState {
	res := make([]regources.ConvertedBalanceState, 0, len(states))

	for _, state := range states {
		res = append(res, regources.ConvertedBalanceState{
			Attributes: regources.ConvertedBalanceStateAttributes{
				InitialAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(amount.MustParseU(state.initialAvailable)),
				},
				ConvertedAmounts: regources.BalanceStateAttributeAmounts{
					Available: regources.Amount(amount.MustParseU(state.convertedAvailable)),
				},
				IsConverted: state.isConverted,
			},
		})
	}

	return res
}

func TestSortConvertedStates(t *testing.T) {
	initialStates := []sortEntry{
		{
			convertedAvailable: "0",
			initialAvailable:   "0",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "0.200000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "3.500000",
			isConverted:        false,
		},
		{
			convertedAvailable: "4.000000",
			initialAvailable:   "5.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "4.000000",
			isConverted:        false,
		},
		{
			convertedAvailable: "2.000000",
			initialAvailable:   "3.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "6.000000",
			initialAvailable:   "7.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "3.200000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "5.000000",
			isConverted:        false,
		},
		{
			convertedAvailable: "5.500000",
			initialAvailable:   "72.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "0",
			isConverted:        true,
		},
	}

	expectedStates := []sortEntry{
		{
			convertedAvailable: "6.000000",
			initialAvailable:   "7.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "5.500000",
			initialAvailable:   "72.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "4.000000",
			initialAvailable:   "5.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "2.000000",
			initialAvailable:   "3.000000",
			isConverted:        true,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "5.000000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "4.000000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "3.500000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "3.200000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "0.200000",
			isConverted:        false,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "0",
			isConverted:        true,
		},
		{
			convertedAvailable: "0",
			initialAvailable:   "0",
			isConverted:        false,
		},
	}

	t.Run("correctly sort the state", func(t *testing.T) {
		sortedStates := SortConvertedStates(fromSortEntries(initialStates))
		assert.Equal(t, sortedStates, fromSortEntries(expectedStates))
	})
}
