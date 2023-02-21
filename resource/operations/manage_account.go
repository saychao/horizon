package operations

type ManageAccount struct {
	Base
	Account              string `json:"account,omitempty"`
	BlockReasonsToAdd    uint32 `json:"block_reasons_to_add,omitempty"`
	BlockReasonsToRemove uint32 `json:"block_reasons_to_remove,omitempty"`
}
