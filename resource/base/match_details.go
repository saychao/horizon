package base

type MatchDetails struct {
	BaseAmount  string `json:"base_amount"`
	QuoteAmount string `json:"quote_amount"`
	FeePair     string `json:"fee_pair"`
	Price       string `json:"price"`
}
