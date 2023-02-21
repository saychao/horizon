package resources

import (
	"fmt"
	"gitlab.com/tokend/regources/generated"
)

// NewOrderBookKey - creates new Key for OrderBook
func NewOrderBookKey(base string, quote string, id uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprintf("%s:%s:%d", base, quote, id),
		Type: regources.ORDER_BOOKS,
	}
}

// NewOrderBook - creates new instance of OrderBook
func NewOrderBook(base string, quote string, id uint64) regources.OrderBook {
	return regources.OrderBook{
		Key: NewOrderBookKey(base, quote, id),
		Relationships: regources.OrderBookRelationships{
			BaseAsset:  NewAssetKey(base).AsRelation(),
			QuoteAsset: NewAssetKey(quote).AsRelation(),
		},
	}
}
