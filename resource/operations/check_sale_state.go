package operations

type CheckSaleState struct {
	Base
	SaleID uint64 `json:"sale_id"`
	Effect string `json:"effect"`
}
