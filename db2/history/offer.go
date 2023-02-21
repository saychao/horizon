package history

import (
	"github.com/saychao/horizon/db2"
	"time"
)

type Offer struct {
	db2.TotalOrderID
	OfferID           int64     `db:"offer_id"`
	OwnerID           string    `db:"owner_id"`
	BaseAsset         string    `db:"base_asset"`
	QuoteAsset        string    `db:"quote_asset"`
	IsBuy             bool      `db:"is_buy"`
	InitialBaseAmount int64     `db:"initial_base_amount"`
	CurrentBaseAmount int64     `db:"current_base_amount"`
	Price             int64     `db:"price"`
	IsCanceled        bool      `db:"is_canceled"`
	CreatedAt         time.Time `db:"created_at"`
}
