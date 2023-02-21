package ingest

import (
	"github.com/saychao/horizon/db2/core"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Load runs queries against `core` to fill in the records of the bundle.
func (lb *LedgerBundle) Load(coreQ core.QInterface) error {

	// Load Header
	err := coreQ.LedgerHeaderBySequence(&lb.Header, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to get ledger header")
	}

	// Load transactions
	err = coreQ.TransactionsByLedger(&lb.Transactions, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to get transactions for ledger")
	}

	err = coreQ.TransactionFeesByLedger(&lb.TransactionFees, lb.Sequence)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction fees for ledger")
	}

	return nil
}
