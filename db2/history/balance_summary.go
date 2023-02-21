package history

import (
	"encoding/json"
	"fmt"
	"time"
)

type BalanceSummary struct {
	BalanceID    string                `db:"balance_id"`
	AmountBefore int64                 `db:"amount_before"`
	Updates      BalanceSummaryUpdates `db:"updates"`
}

type BalanceSummaryUpdate struct {
	UpdatedAt time.Time `json:"updated_at"`
	Amount    int64     `json:"amount"`
}

type BalanceSummaryUpdates []BalanceSummaryUpdate

func (updates *BalanceSummaryUpdates) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &updates)
	default:
		return fmt.Errorf("unsupported Scan from type %T", v)
	}
}
