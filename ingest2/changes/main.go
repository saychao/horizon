// Package changes provides handlers which update state of entities in the DB according to change of the entry in core
package changes

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type ledgerChange struct {
	LedgerSeq       int32
	LedgerCloseTime time.Time
	LedgerChange    xdr.LedgerEntryChange
	Operation       *xdr.Operation
	OperationResult *xdr.OperationResultTr
	OperationIndex  uint32
	OpChanges       xdr.LedgerEntryChanges
	TxIndex         uint32
}

func unixToTime(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}
