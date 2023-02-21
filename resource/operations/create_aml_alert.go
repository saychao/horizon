package operations

type CreateAmlAlert struct {
	Base
	BalanceID string `json:"balance_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
	Reference string `json:"reference"`
}
