package resource

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

// Populate fills out the details
func PopulateTransactionV2(
	row history.Transaction, ledgerChanges []history.LedgerChanges,
) (t regources.Transaction) {
	t = PopulateTransaction(row)
	for _, change := range ledgerChanges {
		ledgerEntryChangeV2 := regources.LedgerEntryChangeV2{
			Effect:    int32(change.Effect),
			EntryType: int32(change.EntryType),
			Payload:   change.Payload,
		}
		t.Changes = append(t.Changes, ledgerEntryChangeV2)
	}
	return t
}
