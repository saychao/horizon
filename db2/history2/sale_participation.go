package history2

// SaleParticipation is a db representations of the sale participation
type SaleParticipation struct {
	ID            uint64 `db:"id"`
	SaleID        uint64 `db:"sale_id"`
	ParticipantID string `db:"participant_id"`
	BaseAmount    string `db:"base_amount"`
	QuoteAmount   string `db:"quote_amount"`
	BaseAsset     string `db:"base_asset"`
	QuoteAsset    string `db:"quote_asset"`
}
