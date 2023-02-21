package participants

import (
	"encoding/json"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

type Match struct {
	BaseAmount  xdr.Int64
	QuoteAmount xdr.Int64
	FeePaid     xdr.Int64
	Price       xdr.Int64
}

func NewMatch(baseAmount, quoteAmount, feePaid, price xdr.Int64) *Match {
	return &Match{
		BaseAmount:  baseAmount,
		QuoteAmount: quoteAmount,
		FeePaid:     feePaid,
		Price:       price,
	}
}

func (m *Match) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"base_amount":  amount.String(int64(m.BaseAmount)),
		"quote_amount": amount.String(int64(m.QuoteAmount)),
		"fee_paid":     amount.String(int64(m.FeePaid)),
		"price":        amount.String(int64(m.Price)),
	})
}

func (m *Match) Add(o *Match) {
	m.BaseAmount += o.BaseAmount
	m.QuoteAmount += o.QuoteAmount
	m.FeePaid += o.FeePaid
}

func (m *Match) CanAdd(o *Match) bool {
	return m.Price == o.Price
}
