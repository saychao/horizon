package ingest2

import (
	core "github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/generator"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type ledgerChangesStorage interface {
	Insert(ledgerChanges []history2.LedgerChanges) error
}

// LedgerChangesSaver - handles each ledger to store sequence of changes occurred for ledger entries
type LedgerChangesSaver struct {
	storage ledgerChangesStorage
}

// NewLedgerChangesHandler - creates new instance of LedgerChangesSaver
func NewLedgerChangesHandler(storage ledgerChangesStorage) *LedgerChangesSaver {
	return &LedgerChangesSaver{
		storage: storage,
	}
}

// Handle - stores ledger changes into db
func (h *LedgerChangesSaver) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	txIDGen := generator.NewIDI32(header.Sequence)
	opIDGen := generator.NewIDI32(header.Sequence)
	for txI := range txs {
		txID := txIDGen.Next()
		ops := txs[txI].ResultMeta.MustOperations()
		for opI := range ops {
			opID := opIDGen.Next()
			err := h.handleOpChanges(txID, opID, ops[opI].Changes)
			if err != nil {
				return errors.Wrap(err, "failed to handle op ledger changes", logan.F{
					"tx_id":      txID,
					"op_id":      opID,
					"ledger_seq": header.Sequence,
				})
			}
		}
	}

	return nil
}

func (h *LedgerChangesSaver) handleOpChanges(txID, opID int64, ledgerChanges []xdr.LedgerEntryChange) error {
	toStore := make([]history2.LedgerChanges, len(ledgerChanges))
	for i, change := range ledgerChanges {
		toStore[i] = history2.LedgerChanges{
			TransactionID: txID,
			OperationID:   opID,
			OrderNumber:   i,
			Effect:        int(change.Type),
			EntryType:     int(change.EntryType()),
			Payload:       ledgerChanges[i],
		}
	}

	err := h.storage.Insert(toStore)
	if err != nil {
		return errors.Wrap(err, "failed to store ledger change")
	}

	return nil
}

// Name - name of the handler
func (h *LedgerChangesSaver) Name() string {
	return "ledger_changes_saver"
}
