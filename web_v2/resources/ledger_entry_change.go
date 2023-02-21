package resources

import (
	"fmt"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func newLedgerEntryID(historyChange history2.LedgerChanges) string {
	// TODO: temporary solution, probably we should store changes with ID's
	return fmt.Sprintf(
		"%d:%d:%d",
		historyChange.TransactionID,
		historyChange.OperationID,
		historyChange.OrderNumber,
	)
}

// NewLedgerEntryChangeKey - creates new key for LedgerEntryChange
func NewLedgerEntryChangeKey(historyChange history2.LedgerChanges) regources.Key {
	return regources.Key{
		ID:   newLedgerEntryID(historyChange),
		Type: regources.LEDGER_ENTRY_CHANGES,
	}
}

// NewLedgerEntryChange - creates new instance of LedgerEntryChange
func NewLedgerEntryChange(historyChange history2.LedgerChanges) (*regources.LedgerEntryChange, error) {
	rawPayload, err := historyChange.Payload.Value()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value from ledger entry change payload")
	}

	strPayload, ok := rawPayload.(string)
	if !ok {
		return nil, errors.Wrap(err, "failed to read ledger entry change payload as a string")
	}

	return &regources.LedgerEntryChange{
		Key: NewLedgerEntryChangeKey(historyChange),
		Attributes: regources.LedgerEntryChangeAttributes{
			Payload:    strPayload,
			ChangeType: xdr.LedgerEntryChangeType(historyChange.Effect),
			EntryType:  xdr.LedgerEntryType(historyChange.EntryType),
		},
	}, nil
}
