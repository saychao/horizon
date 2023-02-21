package history2

import (
	"encoding/json"

	regources "gitlab.com/tokend/regources/generated"
)

type DeferredPayment struct {
	ID                 int64                `db:"id"`
	Amount             regources.Amount     `db:"amount"`
	Details            json.RawMessage      `db:"details"`
	SourceAccount      string               `db:"source_account"`
	SourceBalance      string               `db:"source_balance"`
	DestinationAccount string               `db:"destination_account"`
	State              DeferredPaymentState `db:"state"`
}
