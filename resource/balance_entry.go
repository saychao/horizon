package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

type BalanceEntry struct {
	AccountID string `json:"account_id"`
	BalanceID string `json:"balance_id"`
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Locked    string `json:"locked"`
}

func (r *BalanceEntry) Populate(entry xdr.BalanceEntry) {
	r.AccountID = entry.AccountId.Address()
	r.BalanceID = entry.BalanceId.AsString()
	r.Asset = string(entry.Asset)
	r.Amount = amount.String(int64(entry.Amount))
	r.Locked = amount.String(int64(entry.Locked))
}

type LedgerKeyBalance struct {
	BalanceID string `json:"balance_id"`
}

func (r *LedgerKeyBalance) Populate(entry xdr.LedgerKeyBalance) {
	r.BalanceID = entry.BalanceId.AsString()
}
