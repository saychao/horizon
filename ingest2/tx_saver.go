package ingest2

import (
	"github.com/SafeRE-IT/horizon/db2"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	core "github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/generator"
)

type txStorage interface {
	Insert(txs []history2.Transaction) error
}

// TxSaver - converts core tx into horizon tx and stores them to db
type TxSaver struct {
	storage txStorage
}

// NewTxSaver - creates new instance of TxSaver
func NewTxSaver(storage txStorage) *TxSaver {
	return &TxSaver{
		storage: storage,
	}
}

// Handle - converts tx into history tx and stores them into db
func (h *TxSaver) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	toStore := make([]history2.Transaction, len(txs))
	idGenerator := generator.NewID(uint32(header.Sequence))
	for i, tx := range txs {
		hTx := history2.Transaction{
			TotalOrderID: db2.TotalOrderID{
				ID: idGenerator.Next(),
			},
			Hash:             tx.TransactionHash,
			LedgerSequence:   header.Sequence,
			LedgerCloseTime:  time.Unix(header.CloseTime, 0).UTC(),
			ApplicationOrder: int32(i),
			Account:          tx.Envelope.Tx.SourceAccount.Address(),
			OperationCount:   int32(len(tx.Envelope.Tx.Operations)),
			Envelope:         tx.MustEnvelopeXDR(),
			Result:           tx.MustResultXDR(),
			Meta:             tx.MustResultMetaXDR(),
			ValidAfter:       time.Unix(int64(tx.Envelope.Tx.TimeBounds.MinTime), 0).UTC(),
			ValidBefore:      time.Unix(int64(tx.Envelope.Tx.TimeBounds.MaxTime), 0).UTC(),
		}

		toStore[i] = hTx
	}

	err := h.storage.Insert(toStore)
	if err != nil {
		return errors.Wrap(err, "failed to insert txs", logan.F{
			"ledger_seq": header.Sequence,
		})
	}

	return nil
}

//Name - returns name of the handler
func (h *TxSaver) Name() string {
	return "tx_saver"
}
