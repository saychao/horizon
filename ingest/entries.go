package ingest

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

func (is *Session) operationChanges(changes xdr.LedgerEntryChanges) error {
	for i := range changes {
		var err error

		switch changes[i].Type {
		case xdr.LedgerEntryChangeTypeCreated:
			err = is.operationCreatedEntry(changes[i].Created)
		case xdr.LedgerEntryChangeTypeRemoved:
			err = is.operationDeletedEntry(changes[i].Removed)
		case xdr.LedgerEntryChangeTypeUpdated:
			err = is.operationUpdatedEntry(changes[i].Updated)
		default:
			// nothing to do here
		}
		if err != nil {
			return errors.Wrap(err, "failed to process operation changes")
		}

		err = is.ledgerChanges(i, changes[i])
		if err != nil {
			return errors.Wrap(err, "failed to process ledger changes", logan.F{
				"change": changes[i],
			})
		}
	}

	return nil
}

func (is *Session) ledgerChanges(orderNumber int, change xdr.LedgerEntryChange) error {
	payload, entryType, ok := getLedgerKeyOrEntry(change)
	if !ok {
		return nil
	}

	err := is.Ingestion.LedgerChanges(
		is.Cursor.TransactionID(),
		is.Cursor.OperationID(),
		orderNumber,
		int(change.Type),
		entryType,
		payload)
	if err != nil {
		return errors.Wrap(err, "failed to ingest ledger changes")
	}

	return nil
}

func getLedgerKeyOrEntry(change xdr.LedgerEntryChange) (interface{}, xdr.LedgerEntryType, bool) {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		return change.MustCreated(), change.MustCreated().Data.Type, true
	case xdr.LedgerEntryChangeTypeRemoved:
		return change.MustRemoved(), change.MustRemoved().Type, true
	case xdr.LedgerEntryChangeTypeUpdated:
		return change.MustUpdated(), change.MustUpdated().Data.Type, true
	default:
		return new(interface{}), xdr.LedgerEntryType(1), false
	}
}

func (is *Session) operationCreatedEntry(ledgerEntry *xdr.LedgerEntry) error {
	handler, ok := creationHandlers[ledgerEntry.Data.Type]
	if !ok {
		return nil
	}

	return handler(is, ledgerEntry)
}

func (is *Session) operationDeletedEntry(ledgerKey *xdr.LedgerKey) error {
	handler, ok := deletionHandlers[ledgerKey.Type]
	if !ok {
		return nil
	}

	return handler(is, ledgerKey)
}

func (is *Session) operationUpdatedEntry(ledgerEntry *xdr.LedgerEntry) error {
	handler, ok := updateHandlers[ledgerEntry.Data.Type]
	if !ok {
		return nil
	}

	return handler(is, ledgerEntry)
}

var creationHandlers = map[xdr.LedgerEntryType]func(is *Session, ledgerEntry *xdr.LedgerEntry) error{
	xdr.LedgerEntryTypeBalance:           balanceCreated,
	xdr.LedgerEntryTypeReviewableRequest: reviewableRequestCreate,
	xdr.LedgerEntryTypeSale:              saleCreate,
	xdr.LedgerEntryTypeContract:          contractCreate,
}

var deletionHandlers = map[xdr.LedgerEntryType]func(is *Session, ledgerKey *xdr.LedgerKey) error{
	xdr.LedgerEntryTypeReviewableRequest: reviewableRequestDelete,
	xdr.LedgerEntryTypeContract:          contractDelete,
}

var updateHandlers = map[xdr.LedgerEntryType]func(is *Session, ledgerKey *xdr.LedgerEntry) error{
	xdr.LedgerEntryTypeBalance:           balanceUpdated,
	xdr.LedgerEntryTypeReviewableRequest: reviewableRequestUpdate,
	xdr.LedgerEntryTypeSale:              saleUpdate,
	xdr.LedgerEntryTypeContract:          contractUpdate,
}
