package history2

import (
	"gitlab.com/tokend/go/xdr"
)

// AccountSpecificRule is a row of data from the `account_specific_rules` table
type AccountSpecificRule struct {
	ID        uint64    `db:"id"`
	Address   *string   `db:"address"`
	Forbids   bool      `db:"forbids"`
	EntryType int32     `db:"entry_type"`
	Key       LedgerKey `db:"key"`
}

func NewAccountSpecificRule(entry xdr.AccountSpecificRuleEntry) AccountSpecificRule {
	address := entry.AccountId.Address()
	var addressPtr *string
	if address != "" {
		addressPtr = &address
	}

	return AccountSpecificRule{
		ID:        uint64(entry.Id),
		Key:       LedgerKey(entry.LedgerKey),
		Address:   addressPtr,
		Forbids:   entry.Forbids,
		EntryType: int32(entry.LedgerKey.Type),
	}
}
