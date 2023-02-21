package participants

import (
	"gitlab.com/tokend/go/xdr"
)

type MatchIngestDetails struct {
	BaseAsset  string  `json:"base_asset"`
	QuoteAsset string  `json:"quote_asset"`
	IsBuy      bool    `json:"is_buy"`
	Matches    []Match `json:"matches"`
}

type MatchesDetails struct {
	Source     xdr.AccountId
	Balance    xdr.BalanceId
	BaseAsset  string
	QuoteAsset string
	IsBuy      bool
	Matches    []Match
}

func NewMatchesDetails(source xdr.AccountId, balance xdr.BalanceId, baseAsset, quoteAsset string, isBuy bool) *MatchesDetails {
	return &MatchesDetails{
		Source:     source,
		Balance:    balance,
		BaseAsset:  baseAsset,
		QuoteAsset: quoteAsset,
		IsBuy:      isBuy,
		Matches:    []Match{},
	}
}

func (m *MatchesDetails) ToParticipant() Participant {
	return Participant{m.Source, &m.Balance, []MatchIngestDetails{
		{
			BaseAsset:  m.BaseAsset,
			QuoteAsset: m.QuoteAsset,
			IsBuy:      m.IsBuy,
			Matches:    m.Matches,
		},
	}}
}

func (m *MatchesDetails) Add(match *Match) {
	for i := range m.Matches {
		if m.Matches[i].CanAdd(match) {
			m.Matches[i].Add(match)
			return
		}
	}

	m.Matches = append(m.Matches, *match)
}
