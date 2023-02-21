package generator

import (
	"math"
	"testing"
)

func TestIdGenerator_Next(t *testing.T) {
	inputs := []struct {
		significant    int32
		seq            uint32
		expectedResult int64
	}{
		{
			0, 0, 1,
		},
		{
			0, 1, 2,
		},
		{
			math.MaxInt32, math.MaxUint32 - 1, math.MaxInt64,
		},
	}

	for _, input := range inputs {
		generator := ID{
			Significant: input.significant,
			seq:         input.seq,
		}

		actual := generator.Next()
		if actual != input.expectedResult {
			t.Fatalf("failed to generate id: expected: %d, actual: %d", input.expectedResult, actual)
		}
	}
}
