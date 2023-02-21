// Package results provides an implementation of the txsub.ResultProvider interface
// backed using the SQL databases used by both stellar core and horizon
package txsub

import (
	"bytes"
	"database/sql"
	"encoding/base64"

	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/txsub"
	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/db2/history"
	"golang.org/x/net/context"
)

// ResultsProvider provides transaction submission results by querying the
// connected horizon and stellar core databases.
type ResultsProvider struct {
	Core    core.QInterface
	History history.QInterface
}

// ResultByHash implements txsub.ResultProvider
func (rp *ResultsProvider) ResultByHash(ctx context.Context, hash string) *txsub.Result {
	// query history database
	var hr history.Transaction
	err := rp.History.TransactionByHash(&hr, hash)
	if err == nil {
		return txResultFromHistory(hr)
	}

	if err != sql.ErrNoRows {
		return &txsub.Result{Err: logan.Wrap(err, "Failed to load tx from horizon db")}
	}

	// query core database
	var cr core.Transaction
	err = rp.Core.TransactionByHash(&cr, hash)
	if err == nil {
		return txResultFromCore(cr)
	}

	if err != sql.ErrNoRows {
		return &txsub.Result{Err: logan.Wrap(err, "Failed to load tx from core db")}
	}

	// if no result was found in either db, return nil
	return nil
}

func txResultFromHistory(tx history.Transaction) *txsub.Result {
	return &txsub.Result{
		Hash:           tx.TransactionHash,
		LedgerSequence: tx.LedgerSequence,
		EnvelopeXDR:    tx.TxEnvelope,
		ResultXDR:      tx.TxResult,
		ResultMetaXDR:  tx.TxMeta,
	}
}

func txResultFromCore(tx core.Transaction) *txsub.Result {
	// re-encode result to base64
	var raw bytes.Buffer
	_, err := xdr.Marshal(&raw, tx.Result.Result)

	if err != nil {
		return &txsub.Result{Err: logan.Wrap(err, "Failed to marshal tx result into xdr")}
	}

	trx := base64.StdEncoding.EncodeToString(raw.Bytes())

	// if result is success, send a normal response
	if tx.Result.Result.Result.Code == xdr.TransactionResultCodeTxSuccess {
		return &txsub.Result{
			Hash:           tx.TransactionHash,
			LedgerSequence: tx.LedgerSequence,
			EnvelopeXDR:    tx.EnvelopeXDR(),
			ResultXDR:      trx,
			ResultMetaXDR:  tx.ResultMetaXDR(),
		}
	}

	// if failed, produce a FailedTransactionError
	return &txsub.Result{
		Err:            txsub.NewRejectedTxError(trx),
		Hash:           tx.TransactionHash,
		LedgerSequence: tx.LedgerSequence,
		EnvelopeXDR:    tx.EnvelopeXDR(),
		ResultXDR:      trx,
		ResultMetaXDR:  tx.ResultMetaXDR(),
	}
}
