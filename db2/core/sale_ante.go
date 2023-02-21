package core

type SaleAnte struct {
	SaleID               uint64 `db:"sale_id"`
	ParticipantBalanceID string `db:"participant_balance_id"`
	Amount               uint64 `db:"amount"`
}
