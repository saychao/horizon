package requests

import (
	"fmt"
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

const (
	// IncludeTypeOrderBookBuyEntries - defines if buy entries should be included in the response
	IncludeTypeOrderBookBuyEntries = "buy_entries"
	// IncludeTypeOrderBookSellEntries - defines if sell entries should be included in the response
	IncludeTypeOrderBookSellEntries = "sell_entries"
	// IncludeTypeOrderBookBaseAsset - defines if base asset should be included in the response
	IncludeTypeOrderBookBaseAsset = "base_asset"
	// IncludeTypeOrderBookQuoteAsset - defines if quote asset should be included in the response
	IncludeTypeOrderBookQuoteAsset = "quote_asset"
)

var includeTypeOrderBookAll = map[string]struct{}{
	IncludeTypeOrderBookBuyEntries:  {},
	IncludeTypeOrderBookSellEntries: {},
	IncludeTypeOrderBookBaseAsset:   {},
	IncludeTypeOrderBookQuoteAsset:  {},
}

// GetOrderBook represents params to be specified by user for getOrderBook handler
type GetOrderBook struct {
	*base

	BaseAsset   string
	QuoteAsset  string
	OrderBookID uint64
	MaxEntries  uint64
}

// NewGetOrderBook - returns new instance of GetOrderBook
func NewGetOrderBook(r *http.Request) (*GetOrderBook, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeOrderBookAll,
	})
	if err != nil {
		return nil, err
	}

	baseAsset := b.getString("base")
	if baseAsset == "" {
		return nil, validation.Errors{
			"id": errors.New("id format is: 'base:quote:order_book_id', but base is empty"),
		}
	}
	quoteAsset := b.getString("quote")
	if quoteAsset	 == "" {
		return nil, validation.Errors{
			"id": errors.New("id format is: 'base:quote:order_book_id', but quote is empty"),
		}
	}
	orderBookID, err := b.getUint64("order_book_id")
	if err != nil {
		return nil, validation.Errors{
			"id":   errors.New("id format is: 'base:quote:order_book_id', but order_book_id is invalid"),
		}
	}

	maxEntries, err := b.getUint64("max_entries")
	if err != nil {
		return nil, err
	}
	if maxEntries > maxLimit {
		return nil, validation.Errors{
			"max_entries": errors.New(fmt.Sprintf("The maximum value is %d", maxLimit)),
		}
	}
	if maxEntries == 0 {
		maxEntries = defaultLimit
	}

	request := GetOrderBook{
		base:        b,
		BaseAsset:   baseAsset,
		QuoteAsset:  quoteAsset,
		OrderBookID: orderBookID,
		MaxEntries:  maxEntries,
	}

	return &request, nil
}
