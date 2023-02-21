package ingest2

import core "github.com/SafeRE-IT/horizon/db2/core2"

// LedgerBundle represents a single ledger's worth of novelty created by one
// ledger close
type LedgerBundle struct {
	Header       core.LedgerHeader
	Transactions []core.Transaction
}
