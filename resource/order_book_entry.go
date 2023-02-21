package resource

import (
	"strconv"

	"time"

	"github.com/saychao/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
)

type OrderBookEntry struct {
	PT      string `json:"paging_token"`
	OfferID uint64 `json:"offer_id,omitempty"`
	OwnerID string `json:"owner_id,omitempty"`
	OfferData
}

func (o *OrderBookEntry) Populate(s *core.OrderBookEntry, baseAsset, quoteAsset string, isBuy bool) {
	o.OfferData.Populate(s, baseAsset, quoteAsset, isBuy)
	o.PT = o.PagingToken()

}

// PagingToken implementation for hal.Pageable
func (o *OrderBookEntry) PagingToken() string {
	return strconv.FormatInt(int64(amount.MustParse(o.Price)), 10)
}

type OfferData struct {
	BaseAssetCode  string    `json:"base_asset_code"`
	QuoteAssetCode string    `json:"quote_asset_code"`
	IsBuy          bool      `json:"is_buy"`
	BaseAmount     string    `json:"base_amount"`
	QuoteAmount    string    `json:"quote_amount"`
	Price          string    `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
}

func (o *OfferData) Populate(s *core.OrderBookEntry, baseAsset, quoteAsset string, isBuy bool) {
	o.BaseAssetCode = baseAsset
	o.QuoteAssetCode = quoteAsset
	o.IsBuy = isBuy
	o.BaseAmount = amount.String(s.BaseAmount)
	o.QuoteAmount = amount.String(s.QuoteAmount)
	o.Price = amount.String(s.Price)
	o.CreatedAt = time.Unix(s.CreatedAt, 0).UTC()
}
