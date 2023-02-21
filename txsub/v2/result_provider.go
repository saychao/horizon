package txsub

import (
	"bytes"
	"encoding/base64"

	"github.com/SafeRE-IT/horizon/ingest2/generator"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
)

// ResultsProvider provides transaction submission results by querying the
// connected horizon and stellar core databases.
type ResultsProvider struct {
	Core    core2.TransactionQ
	History history2.TransactionsQ
}

func (rp *ResultsProvider) FromHistory(hash string) (*Result, error) {
	hr, err := rp.History.GetByHash(hash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get history transaction by hash")
	}
	if hr != nil {
		return txResultFromHistory(*hr)
	}

	return nil, nil
}

func (rp *ResultsProvider) FromCore(hash string) (*Result, error) {
	cr, err := rp.Core.GetByHash(hash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get core transaction by hash")
	}
	if cr != nil {
		return txResultFromCore(*cr)
	}

	return nil, nil
}

func txResultFromHistory(tx history2.Transaction) (*Result, error) {
	return &Result{
		TransactionID:  tx.ID,
		Hash:           tx.Hash,
		LedgerSequence: tx.LedgerSequence,
		EnvelopeXDR:    tx.Envelope,
		ResultXDR:      tx.Result,
		ResultMetaXDR:  tx.Meta,
	}, nil
}

func txResultFromCore(tx core2.Transaction) (*Result, error) {
	// re-encode result to base64
	var raw bytes.Buffer
	_, err := xdr.Marshal(&raw, tx.Result.Result)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal tx result into xdr")
	}

	trx := base64.StdEncoding.EncodeToString(raw.Bytes())

	// if result is success, send a normal response
	if tx.Result.Result.Result.Code == xdr.TransactionResultCodeTxSuccess {
		return &Result{
			TransactionID:  generator.MakeIDUint32(tx.LedgerSequence, uint32(tx.Index)),
			Hash:           tx.TransactionHash,
			LedgerSequence: tx.LedgerSequence,
			EnvelopeXDR:    tx.MustEnvelopeXDR(),
			ResultXDR:      trx,
			ResultMetaXDR:  tx.MustResultMetaXDR(),
		}, nil
	}
	// if failed, produce a FailedTransactionError
	return nil, NewRejectedTxError(trx)

}
