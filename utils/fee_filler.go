package utils

import (
	"math"
	"sort"

	"github.com/SafeRE-IT/horizon/db2/core"
)

type sortedFees []core.FeeEntry

func (s sortedFees) Len() int {
	return len(s)
}

func (s sortedFees) Less(i, j int) bool {
	return s[i].LowerBound < s[j].LowerBound
}

func (s sortedFees) Swap(i, j int) {
	data := s[i]
	s[i] = s[j]
	s[j] = data
}

func (s sortedFees) Add(entry core.FeeEntry) sortedFees {
	result := append(s, entry)
	sort.Sort(result)
	return result
}

func FillFeeGaps(rawFees []core.FeeEntry, zeroFee core.FeeEntry) []core.FeeEntry {
	if len(rawFees) == 0 {
		return rawFees
	}

	fees := sortedFees(rawFees)
	sort.Sort(fees)

	// check lower bound
	// no need to add [0,0] fee
	if fees[0].LowerBound != 0 && fees[0].LowerBound != 1 {
		fees = fees.Add(getNewZeroFee(zeroFee, 0, fees[0].LowerBound-1))
	}

	// check upper bound
	if fees[fees.Len()-1].UpperBound != math.MaxInt64 {
		fees = fees.Add(getNewZeroFee(zeroFee, fees[fees.Len()-1].UpperBound+1, math.MaxInt64))
	}

	for i := 0; i < fees.Len()-1; i++ {
		if fees[i].UpperBound == math.MaxInt64 {
			break
		}

		expectedLowerBoundForNextFee := fees[i].UpperBound + 1
		if expectedLowerBoundForNextFee != fees[i+1].LowerBound {
			fees = fees.Add(getNewZeroFee(zeroFee, expectedLowerBoundForNextFee, fees[i+1].LowerBound-1))
		}
	}

	return fees
}

func getNewZeroFee(zeroFee core.FeeEntry, lowerBound, upperBound int64) core.FeeEntry {
	zeroFee.LowerBound = lowerBound
	zeroFee.UpperBound = upperBound
	return zeroFee
}
