package history2

import (
	"time"

	"gitlab.com/tokend/regources/generated"
)

// CreateSwap - represents instance of voting campaign
type Swap struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	LockTime  time.Time `db:"lock_time"`

	SourceAccount      string `db:"source_account"`
	DestinationAccount string `db:"destination_account"`

	SourceBalance      string `db:"source_balance"`
	DestinationBalance string `db:"destination_balance"`

	Secret     *string `db:"secret"`
	SecretHash string  `db:"secret_hash"`

	Amount uint64 `db:"amount"`
	Asset  string `db:"asset"`

	SourceFixedFee        uint64 `db:"source_fixed_fee"`
	SourcePercentFee      uint64 `db:"source_percent_fee"`
	DestinationFixedFee   uint64 `db:"destination_fixed_fee"`
	DestinationPercentFee uint64 `db:"destination_percent_fee"`

	Details regources.Details   `db:"details"`
	State   regources.SwapState `db:"state"`
}
