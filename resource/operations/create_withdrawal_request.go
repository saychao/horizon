package operations

type CreateWithdrawalRequest struct {
	Base
	Amount          string                 `json:"amount"`
	Balance         string                 `json:"balance"`
	FeeFixed        string                 `json:"fee_fixed"`
	FeePercent      string                 `json:"fee_percent"`
	ExternalDetails map[string]interface{} `json:"external_details"`
	DestAsset       string                 `json:"dest_asset"`
	DestAmount      string                 `json:"dest_amount"`
	AllTasks        *uint32                `json:"all_tasks"`
}
