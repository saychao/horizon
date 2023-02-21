package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/SafeRE-IT/horizon/db2/core"
)

var Tests = []struct {
	coreFees []core.FeeEntry
	zeroFee  core.FeeEntry
	expected []core.FeeEntry
}{
	{nil, core.FeeEntry{}, nil},
	// gap at the begging, mid, end
	{[]core.FeeEntry{
		{LowerBound: 402, UpperBound: 500, Percent: 4},
		{LowerBound: 100, UpperBound: 100, Percent: 2},
		{LowerBound: 101, UpperBound: 200, Percent: 3},
		{LowerBound: 300, UpperBound: 400, Percent: 4},
	},
		core.FeeEntry{},
		[]core.FeeEntry{
			{LowerBound: 0, UpperBound: 99},
			{LowerBound: 100, UpperBound: 100, Percent: 2},
			{LowerBound: 101, UpperBound: 200, Percent: 3},
			{LowerBound: 201, UpperBound: 299, Percent: 0},
			{LowerBound: 300, UpperBound: 400, Percent: 4},
			{LowerBound: 401, UpperBound: 401, Percent: 0},
			{LowerBound: 402, UpperBound: 500, Percent: 4},
			{LowerBound: 501, UpperBound: math.MaxInt64, Percent: 0},
		}},
}

func TestParse(t *testing.T) {
	for _, v := range Tests {
		actualResult := FillFeeGaps(v.coreFees, v.zeroFee)
		require.Equal(t, len(v.expected), len(actualResult))
		for i := range v.expected {
			require.Equal(t, v.expected[i], actualResult[i])
		}
	}
}
