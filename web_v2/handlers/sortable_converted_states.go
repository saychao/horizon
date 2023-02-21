package handlers

import (
	"gitlab.com/tokend/regources/generated"
	"sort"
)

func SortConvertedStates(states []regources.ConvertedBalanceState) []regources.ConvertedBalanceState {
	sort.Sort(sortableConvertedStates(states))

	return states
}

type sortableConvertedStates []regources.ConvertedBalanceState

func (s sortableConvertedStates) Len() int {
	return len(s)
}

func (s sortableConvertedStates) Less(i, j int) bool {
	a := s[i].Attributes
	b := s[j].Attributes

	if a.IsConverted && b.IsConverted {
		return a.ConvertedAmounts.Available > b.ConvertedAmounts.Available
	}

	if a.IsConverted && !b.IsConverted {
		return a.ConvertedAmounts.Available > 0 || b.InitialAmounts.Available == 0
	}

	if !a.IsConverted && b.IsConverted {
		return a.InitialAmounts.Available > 0 && b.InitialAmounts.Available == 0
	}

	return a.InitialAmounts.Available > b.InitialAmounts.Available
}

func (s sortableConvertedStates) Swap(i, j int) {
	data := s[i]
	s[i] = s[j]
	s[j] = data
}
