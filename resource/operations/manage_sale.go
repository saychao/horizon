package operations

type ManageSale struct {
	Base
	SaleID    uint64 `json:"sale_id"`
	Action    string `json:"action"`
	Fulfilled bool   `json:"fulfilled"`
}
