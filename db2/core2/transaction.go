package core2

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// Transaction is row of data from the `txhistory` table from stellar-core
type Transaction struct {
	TransactionHash string                    `db:"txid"`
	LedgerSequence  int32                     `db:"ledgerseq"`
	Index           int32                     `db:"txindex"`
	Envelope        xdr.TransactionEnvelope   `db:"txbody"`
	Result          xdr.TransactionResultPair `db:"txresult"`
	ResultMeta      xdr.TransactionMeta       `db:"txmeta"`
}

// IsSuccessful returns true when the transaction was successful.
func (tx *Transaction) IsSuccessful() bool {
	return tx.Result.Result.Result.Code == xdr.TransactionResultCodeTxSuccess
}

// MustEnvelopeXDR returns the XDR encoded envelope for this transaction
func (tx *Transaction) MustEnvelopeXDR() string {
	out, err := xdr.MarshalBase64(tx.Envelope)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal tx", logan.F{
			"tx_hash": tx.TransactionHash,
		}))
	}
	return out
}

// MustResultXDR returns the XDR encoded result for this transaction
func (tx *Transaction) MustResultXDR() string {
	out, err := xdr.MarshalBase64(tx.Result.Result)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal tx result", logan.F{
			"tx_hash": tx.TransactionHash,
		}))
	}
	return out
}

// MustResultMetaXDR returns the XDR encoded result meta for this transaction
func (tx *Transaction) MustResultMetaXDR() string {
	out, err := xdr.MarshalBase64(tx.ResultMeta)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal tx result meta", logan.F{
			"tx_hash": tx.TransactionHash,
		}))
	}
	return out
}
