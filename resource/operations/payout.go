package operations

type Payout struct {
	Base
	Asset                string `json:"asset"`
	SourceBalanceID      string `json:"source_balance_id"`
	MaxPayoutAmount      string `json:"max_payout_amount"`
	ActualPayoutAmount   string `json:"actual_payout_amount"`
	MinPayoutAmount      string `json:"min_payout_amount"`
	MinAssetHolderAmount string `json:"min_asset_holder_amount"`
	FixedFee             string `json:"fixed_fee"`
	PercentFee           string `json:"percent_fee"`
}
