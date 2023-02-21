package core2

// Offer - db representation of offer
type Offer struct {
	OwnerID        string `db:"owner_id"`
	OfferID        uint64 `db:"offer_id"`
	OrderBookID    uint64 `db:"order_book_id"`
	Fee            uint64 `db:"fee"`
	BaseAssetCode  string `db:"base_asset_code"`
	QuoteAssetCode string `db:"quote_asset_code"`
	BaseBalanceID  string `db:"base_balance_id"`
	QuoteBalanceID string `db:"quote_balance_id"`
	BaseAmount     uint64 `db:"base_amount"`
	QuoteAmount    uint64 `db:"quote_amount"`
	IsBuy          bool   `db:"is_buy"`
	CreatedAt      int64  `db:"created_at"`
	Price          int64  `db:"price"`

	BaseAsset  *Asset `db:"base_assets"`
	QuoteAsset *Asset `db:"quote_assets"`
}
