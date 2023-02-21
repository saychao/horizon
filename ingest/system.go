package ingest

import (
	err2 "github.com/pkg/errors"
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/errors"
	"github.com/saychao/horizon/ledger"
	"github.com/saychao/horizon/log"
)

// Tick triggers the ingestion system to ingest any new ledger data, provided
// that there currently is not an import session in progress.
func (i *System) Tick() *Session {
	i.lock.Lock()
	if i.current != nil {
		log.Info("ingest: already in progress")
		i.lock.Unlock()
		return nil
	}

	is := i.newTickSession()
	i.current = is
	i.lock.Unlock()

	i.runOnce()
	return is
}

// newTickSession creates an unverified new ingestion session that reflects the
// current cached ledger state.
func (i *System) newTickSession() *Session {
	var (
		start int32
		ls    = ledger.CurrentState()
	)

	if ls.History.Latest == 0 {
		start = ls.Core.OldestOnStart
	} else {
		start = ls.History.Latest + 1
	}

	end := ls.Core.Latest

	return NewSession(ls.History.OldestOnStart > 1 || ls.Core.OldestOnStart > 1, start, end, i)
}

// run causes the importer to check stellar-core to see if we can import new
// data.
func (i *System) runOnce() {
	defer func() {
		if rec := recover(); rec != nil {
			err := errors.FromPanic(rec)
			log.WithField("service", "ingest_v1").WithError(err).WithStack(err).Error("Import session paniced")
		}
	}()

	ls := ledger.CurrentState()

	// 1. stash a copy of the current ingestion session (assigned from the tick)
	// 2. output "initial ingestion" message if the
	// 3. import until none available

	// 1.
	i.lock.Lock()
	is := i.current
	i.lock.Unlock()

	defer func() {
		i.lock.Lock()
		i.current = nil
		i.lock.Unlock()
	}()

	if is == nil {
		log.Warn("ingest: runOnce ran with a nil current session")
		return
	}

	if is.Cursor.FirstLedger > is.Cursor.LastLedger {
		return
	}

	// 2.
	if ls.History.Latest == 0 {
		log.Infof(
			"history db is empty, starting ingestion from ledger %d",
			is.Cursor.FirstLedger,
		)
	}

	if is.Cursor.FirstLedger != ls.Core.OldestOnStart {
		err := i.validateLedgerChain(is.Cursor.FirstLedger)
		if err != nil {
			log.
				WithField("start", is.Cursor.FirstLedger).
				Errorf("ledger gap detected (possible db corruption): %s", err)
			return
		}
	}

	// 3.
	is.Run()

	return
}

// validateLedgerChain helps to ensure the chain of ledger entries is contiguous
// within horizon.  It ensures the ledger at `seq` is a child of `seq - 1`.
func (i *System) validateLedgerChain(seq int32) error {
	var (
		cur  core.LedgerHeader
		prev core.LedgerHeader
	)

	q := core.NewQ(i.CoreDB)

	err := q.LedgerHeaderBySequence(&cur, seq)
	if err != nil {
		return err2.Wrap(err, "validateLedgerChain: failed to load cur ledger")
	}

	err = q.LedgerHeaderBySequence(&prev, seq-1)
	if err != nil {
		return err2.Wrap(err, "validateLedgerChain: failed to load prev ledger")
	}

	if cur.PrevHash != prev.LedgerHash {
		return err2.New("cur and prev ledger hashes don't match")
	}

	return nil
}
