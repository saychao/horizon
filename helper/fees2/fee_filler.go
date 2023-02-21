package fees2

import (
	"math"
	"sort"
)

type sortedFees []FeeWrapper

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

func (s sortedFees) Add(entry FeeWrapper) sortedFees {
	result := append(s, entry)
	sort.Sort(result)
	return result
}

func fillFeeGaps(primaryFees []FeeWrapper, secondaryFee FeeWrapper) []FeeWrapper {
	if len(primaryFees) == 0 {
		return []FeeWrapper{secondaryFee}
	}
	fees := sortedFees(primaryFees)
	sort.Sort(fees)

	// check lower bound
	if fees[0].LowerBound != 0 && fees[0].LowerBound != 1 {
		fee, ok := fillGap(0, fees[0].LowerBound-1, secondaryFee)
		if ok {
			fees = fees.Add(fee)
		}

	}

	// check upper bound
	if fees[fees.Len()-1].UpperBound != math.MaxInt64 {
		fee, ok := fillGap(fees[fees.Len()-1].UpperBound+1, math.MaxInt64, secondaryFee)
		if ok {
			fees = fees.Add(fee)
		}
	}

	for i := 0; i < fees.Len()-1; i++ {
		if fees[i].LowerBound > secondaryFee.UpperBound {
			break
		}

		expectedLowerBoundForNextFee := fees[i].UpperBound + 1
		if expectedLowerBoundForNextFee == fees[i+1].LowerBound {
			continue
		}

		fee, ok := fillGap(expectedLowerBoundForNextFee, fees[i+1].LowerBound-1, secondaryFee)
		if !ok {
			continue
		}

		fees = fees.Add(fee)
	}

	return fees
}

func fillGap(lowerBound, upperBound int64, fee FeeWrapper) (filler FeeWrapper, ok bool) {
	if ok, l, u := overlap(lowerBound, upperBound, fee.LowerBound, fee.UpperBound); ok {
		fee.LowerBound = l
		fee.UpperBound = u
		return fee, true
	}
	return fee, false
}

func overlap(a, b, c, d int64) (ok bool, from int64, to int64) {
	if a > b || c > d {
		return
	}
	if b-c >= 0 && d-a >= 0 {
		return true, max(a, c), min(b, d)
	}

	return false, 0, 0
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
