package resources

import (
	"github.com/saychao/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewTransactionKey - creates new key for Transaction
func NewTransactionKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, regources.TRANSACTIONS)
}

// NewTransaction - creates new instance of Transaction
func NewTransaction(historyTx history2.Transaction) regources.Transaction {
	return regources.Transaction{
		Key: NewTransactionKey(historyTx.ID),
		Attributes: regources.TransactionAttributes{
			Hash:           historyTx.Hash,
			LedgerSequence: historyTx.LedgerSequence,
			CreatedAt:      historyTx.LedgerCloseTime,
			EnvelopeXdr:    historyTx.Envelope,
			ResultXdr:      historyTx.Result,
			ResultMetaXdr:  historyTx.Meta,
		},
	}
}
