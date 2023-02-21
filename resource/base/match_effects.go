package base

type MatchEffects struct {
	BaseAsset  string         `json:"base_asset"`
	QuoteAsset string         `json:"quote_asset"`
	IsBuy      bool           `json:"is_buy"`
	Matches    []MatchDetails `json:"matches"`
}
