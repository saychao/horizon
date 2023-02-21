package fees2

import (
	"testing"

	"github.com/magiconair/properties/assert"
	core "github.com/SafeRE-IT/horizon/db2/core2"
)

func TestFillFeeGaps(t *testing.T) {
	t.Run("one gap", func(t *testing.T) {
		primaryFees := []FeeWrapper{
			{
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		secondaryFee := FeeWrapper{
			Fee: core.Fee{
				LowerBound: 2,
				UpperBound: 20,
				Percent:    2,
			},
		}

		expected := []FeeWrapper{
			{
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 6,
					UpperBound: 9,
					Percent:    2,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
			{
				Fee: core.Fee{
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
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 20,
					UpperBound: 22,
				},
			},
		}
		secondaryFees := []FeeWrapper{
			{
				Fee: core.Fee{
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 21,
					UpperBound: 25,
					Percent:    3,
				},
			},
		}

		expected := []FeeWrapper{
			{
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 6,
					UpperBound: 9,
					Percent:    2,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 16,
					UpperBound: 19,
					Percent:    2,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 20,
					UpperBound: 22,
				},
			},
			{
				Fee: core.Fee{
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
				Fee: core.Fee{
					LowerBound: 2,
					UpperBound: 20,
					Percent:    2,
				},
			},
		}

		expected := []FeeWrapper{
			{

				Fee: core.Fee{
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
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		var secondaryFees []FeeWrapper
		expected := []FeeWrapper{
			{
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
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
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}
		secondaryFee := FeeWrapper{
			Fee: core.Fee{
				LowerBound: 0,
				UpperBound: 5,
				Percent:    2,
			},
		}

		expected := []FeeWrapper{
			{
				Fee: core.Fee{
					LowerBound: 0,
					UpperBound: 5,
				},
			},
			{
				Fee: core.Fee{
					LowerBound: 10,
					UpperBound: 15,
				},
			},
		}

		got := fillFeeGaps(primaryFees, secondaryFee)

		assert.Equal(t, got, expected)
	})
}
