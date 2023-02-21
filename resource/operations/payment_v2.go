package operations

import (
	"gitlab.com/tokend/regources"
)

type PaymentV2 struct {
	Base
	PaymentID           uint64              `json:"payment_id"`
	From                string              `json:"from,omitempty"`
	To                  string              `json:"to,omitempty"`
	FromBalance         string              `json:"from_balance,omitempty"`
	ToBalance           string              `json:"to_balance,omitempty"`
	Amount              regources.Amount    `json:"amount"`
	Asset               string              `json:"asset"`
	SourceFeeData       regources.FeeDataV2 `json:"source_fee_data"`
	DestinationFeeData  regources.FeeDataV2 `json:"destination_fee_data"`
	SourcePaysForDest   bool                `json:"source_pays_for_dest"`
	Subject             string              `json:"subject"`
	Reference           string              `json:"reference"`
	SourceSentUniversal string              `json:"source_sent_universal"`
}
