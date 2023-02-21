package fees

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/saychao/horizon/db2/core"
)

func TestFillFeeGaps(t *testing.T) {
	t.Run("one gap", func(t *testing.T) {
		primaryFees := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		secondaryFee := FeeWrapper{
			FeeEntry: core.FeeEntry{
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
		}

		expected := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 6,
					UpperBound: 9,
					Percent:    2,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 16,
					UpperBound: 20,
					Percent:    2,
				},
			},
		}
		got := fillFeeGaps(primaryFees, secondaryFee)

		assert.Equal(t, got, expected)
	})

	t.Run("two gaps", func(t *testing.T) {
		primaryFees := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 20,
					UpperBound: 22,
				},
			},
		}
		secondaryFees := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 21,
					UpperBound: 25,
					Percent:    3,
				},
			},
		}

		expected := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 6,
					UpperBound: 9,
					Percent:    2,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 16,
					UpperBound: 19,
					Percent:    2,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 20,
					UpperBound: 22,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 23,
					UpperBound: 25,
					Percent:    3,
				},
			},
		}

		got := primaryFees
		for _, fee := range secondaryFees {
			got = fillFeeGaps(got, fee)
		}
		assert.Equal(t, got, expected)
	})

	t.Run("no primary fees", func(t *testing.T) {
		var primaryFees []FeeWrapper
		secondaryFees := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
				},
			},
		}

		expected := []FeeWrapper{
			{

				FeeEntry: core.FeeEntry{
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
				},
			},
		}

		got := primaryFees
		for _, fee := range secondaryFees {
			got = fillFeeGaps(got, fee)
		}
		assert.Equal(t, got, expected)
	})

	t.Run("no secondary fees", func(t *testing.T) {
		primaryFees := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		var secondaryFees []FeeWrapper
		expected := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}

		got := primaryFees
		for _, fee := range secondaryFees {
			got = fillFeeGaps(got, fee)
		}
		assert.Equal(t, got, expected)
	})

	t.Run("no overlap", func(t *testing.T) {
		primaryFees := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		secondaryFee := FeeWrapper{
			FeeEntry: core.FeeEntry{
				LowerBound: 0,
				UpperBound: 5,
				Percent:    2,
			},
		}

		expected := []FeeWrapper{
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				FeeEntry: core.FeeEntry{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}

		got := fillFeeGaps(primaryFees, secondaryFee)

		assert.Equal(t, got, expected)
	})
}
