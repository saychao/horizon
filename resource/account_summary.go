package resource

import (
	"time"

	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
)

type AccountSummary struct {
	BalanceSummary map[string]BalanceSummary `json:"balance_summary"`
}

func (as *AccountSummary) Populate(records []history.BalanceSummary) {
	as.BalanceSummary = map[string]BalanceSummary{}
	for _, record := range records {
		bs := BalanceSummary{
			AmountBefore: amount.String(record.AmountBefore),
		}
		bs.Updates = make([]BalanceUpdate, 0, len(record.Updates))
		for _, update := range record.Updates {
			bs.Updates = append(bs.Updates, BalanceUpdate{
				Amount:    amount.String(update.Amount),
				Timestamp: update.UpdatedAt,
			})
		}
		as.BalanceSummary[record.BalanceID] = bs
	}
}

type BalanceSummary struct {
	AmountBefore string          `json:"amount_before"`
	Updates      []BalanceUpdate `json:"updates"`
}

type BalanceUpdate struct {
	Amount    string    `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
