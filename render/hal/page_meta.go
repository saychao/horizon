package hal

import (
	"time"

	"gitlab.com/tokend/regources"
)

type PageMeta struct {
	LatestLedger *LatestLedgerMeta        `json:"latest_ledger,omitempty"`
	Count        *regources.RequestsCount `json:"count,omitempty"`
}

type LatestLedgerMeta struct {
	Sequence int32     `json:"sequence"`
	ClosedAt time.Time `json:"closed_at"`
}
