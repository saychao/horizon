package history

import "time"

type Trades struct {
	ID          int64     `db:"id"`
	OrderBookID uint64    `db:"order_book_id"`
	BaseAsset   string    `db:"base_asset"`
	QuoteAsset  string    `db:"quote_asset"`
	BaseAmount  int64     `db:"base_amount"`
	QuoteAmount int64     `db:"quote_amount"`
	Price       int64     `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
}
