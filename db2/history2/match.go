package history2

import (
	"gitlab.com/tokend/go/xdr"
	"time"
)

// Match is a row of data from the `matches` table
type Match struct {
	ID          int64     `db:"id"`
	OperationID int64     `db:"operation_id"`
	OfferID     uint64    `db:"offer_id"`
	BaseAmount  uint64    `db:"base_amount"`
	QuoteAmount uint64    `db:"quote_amount"`
	BaseAsset   string    `db:"base_asset"`
	QuoteAsset  string    `db:"quote_asset"`
	Price       uint64    `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
}

// NewMatch returns new instance of Match
func NewMatch(matchID int64, operationID int64, base, quote xdr.AssetCode, createdAt time.Time, atom xdr.ClaimOfferAtom) Match {
	return Match{
		ID:          matchID,
		OperationID: operationID,
		OfferID:     uint64(atom.OfferId),
		BaseAmount:  uint64(atom.BaseAmount),
		QuoteAmount: uint64(atom.QuoteAmount),
		BaseAsset:   string(base),
		QuoteAsset:  string(quote),
		Price:       uint64(atom.CurrentPrice),
		CreatedAt:   createdAt,
	}
}
