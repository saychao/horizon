package operations

type CreateIssuanceRequest struct {
	Base
	Reference       string                 `json:"reference"`
	Amount          string                 `json:"amount"`
	Asset           string                 `json:"asset"`
	FeeFixed        string                 `json:"fee_fixed"`
	FeePercent      string                 `json:"fee_percent"`
	ExternalDetails map[string]interface{} `json:"external_details"`
	AllTasks        *uint32                `json:"all_tasks"`
}
