package ingest2

import (
	"context"
	"time"

	core "github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ledger"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const (
	waitForLedgerPeriod    = time.Second
	minErrorRecoveryPeriod = time.Second * 2
	maxErrorRecoveryPeriod = time.Minute * 10
)

// historyLedgerProvider - specifies methods required to get ledger from history db
type historyLedgerProvider interface {
	GetBySequence(seq int32) (*history2.Ledger, error)
}

// txProvider - specifies methods required to get data from core db needed for ingest
type txProvider interface {
	// GetByLedgerRange - returns range of transactions applied in ledgers with sequence in [fromSeq,toSeq]
	GetByLedgerRange(fromSeq int32, toSeq int32) ([]core.Transaction, error)
	// GetBySequence - returns ledger header by it's seq
	GetBySequence(seq int32) (*core.LedgerHeader, error)
	// GetBySequenceRange - returns range of ledger headers with sequence in [fromSeq,toSeq]
	GetBySequenceRange(fromSeq int32, toSeq int32) ([]core.LedgerHeader, error)
}

// Producer - worker which is responsible for loading sequence of ledgers and transactions and sending them
// to provided channel.
type Producer struct {
	txProvider      txProvider
	hLedgerProvider historyLedgerProvider
	log             *logan.Entry

	data           chan LedgerBundle
	currentLedger  int32
	batchSize      int32
	getLedgerState func() ledger.SystemState
}

// NewProducer - creates new instance of ingest data produce
func NewProducer(txProvider txProvider, hLedgerProvider historyLedgerProvider, log *logan.Entry, batchSize int32,
	getLedgerState func() ledger.SystemState) *Producer {

	return &Producer{
		log:             log,
		txProvider:      txProvider,
		hLedgerProvider: hLedgerProvider,
		batchSize:       batchSize,
		getLedgerState:  getLedgerState,
	}
}

// Start - starts the Producer in new goroutine. Panics on initialization step if already started.
func (l *Producer) Start(ctx context.Context) chan LedgerBundle {
	if l.data != nil {
		l.log.Panic("Already started")
	}

	l.data = make(chan LedgerBundle, l.batchSize)
	ledgerState := l.getLedgerState()
	var err error
	l.currentLedger, err = l.getLedgerSeqToStartFrom(ledgerState)
	if err != nil {
		l.log.WithError(err).Panic("Failed to figure out ledger sequence from which to start")
	}

	err = l.ensureChainIsConsistent(l.currentLedger)
	if err != nil {
		l.log.WithError(err).Panic("Failed to ensure consistency of core and horizon chains")
	}
	go running.WithBackOff(ctx, l.log, "ingest_producer", l.runOnce, waitForLedgerPeriod,
		minErrorRecoveryPeriod, maxErrorRecoveryPeriod)

	return l.data
}

func (l *Producer) runOnce(ctx context.Context) error {
	coreLatestVersion := l.getLedgerState().Core.Latest
	if coreLatestVersion < l.currentLedger {
		return errors.From(errors.New("Unexpected state: latest ledger in core is below current ledger in horizon"), logan.F{
			"core_latest":    coreLatestVersion,
			"horizon_latest": l.currentLedger,
		})
	}

	if coreLatestVersion == l.currentLedger {
		return nil
	}

	fromSeq := l.currentLedger + 1
	toSeq := minOf(coreLatestVersion, l.currentLedger+l.batchSize)
	batch, err := l.getBatch(fromSeq, toSeq)
	if err != nil {
		return errors.Wrap(err, "failed to load batch", logan.F{"fromSeq": fromSeq, "toSeq": toSeq})
	}

	for _, bundle := range batch {
		l.data <- bundle
	}
	l.currentLedger = toSeq

	return nil
}

func minOf(values ...int32) int32 {
	if len(values) < 1 {
		panic("Incorrect number of arguments")
	}

	min := values[0]

	for _, value := range values {
		if value < min {
			min = value
		}
	}

	return min
}

func (l *Producer) getBatch(fromSeq int32, toSeq int32) ([]LedgerBundle, error) {
	headers, err := l.txProvider.GetBySequenceRange(fromSeq, toSeq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger headers")
	}

	txs, err := l.txProvider.GetByLedgerRange(fromSeq, toSeq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load ledger transactions")
	}

	result := make([]LedgerBundle, 0, len(headers))
	for _, header := range headers {
		result = append(result, LedgerBundle{
			Header:       header,
			Transactions: make([]core.Transaction, 0, header.Data.MaxTxSetSize),
		})
	}

	for _, tx := range txs {
		if tx.IsSuccessful() {
			bundle := result[tx.LedgerSequence-fromSeq]
			bundle.Transactions = append(bundle.Transactions, tx)
			result[tx.LedgerSequence-fromSeq] = bundle
		}
	}

	return result, nil
}

// ensureChainIsConsistent - ensures that we'll end up with consistent chain of blocks in both core and horizon.
func (l *Producer) ensureChainIsConsistent(ledgerSeqToIngest int32) error {
	// try to load previously ingested ledger
	prevLedger, err := l.hLedgerProvider.GetBySequence(ledgerSeqToIngest - 1)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger by sequence", logan.F{
			"ledger_seq": ledgerSeqToIngest - 1,
		})
	}

	// horizon db is empty, so we have the same chain
	if prevLedger == nil {
		var coreLedger *core.LedgerHeader
		coreLedger, err = l.txProvider.GetBySequence(2)
		if err != nil {
			return errors.Wrap(err, "failed to load ledger from core DB")
		}

		if coreLedger == nil {
			return errors.Wrap(err, "Core does not have full history. Unfortunately this version of horizon does"+
				" not support partial history")
		}
		return nil
	}

	currentLedger, err := l.txProvider.GetBySequence(ledgerSeqToIngest)
	if err != nil {
		return errors.Wrap(err, "failed to load ledger by seq from core", logan.F{
			"ledger_seq": ledgerSeqToIngest,
		})
	}

	if currentLedger == nil {
		return errors.From(errors.New("failed to load from core ledger to ingest - horizon is ahead of core or core is dead"), logan.F{
			"ledger_to_ingest": ledgerSeqToIngest,
		})
	}

	if currentLedger.PrevHash != prevLedger.Hash {
		return errors.From(errors.New("chain in horizon does not match chain in core"), logan.F{
			"current_ledger_prev_hash": currentLedger.PrevHash,
			"prev_ledger_hash":         prevLedger.Hash,
			"prev_ledger_seq":          prevLedger.Sequence,
			"current_ledger_seq":       currentLedger.Sequence,
		})
	}

	return nil

}

func (l *Producer) getLedgerSeqToStartFrom(ledgerState ledger.SystemState) (int32, error) {
	if ledgerState.Core.Latest < ledgerState.History2.Latest {
		return 0, errors.From(errors.New("horizon is ahead of core"), logan.F{
			"core_latest":     ledgerState.Core.Latest,
			"history2_latest": ledgerState.History2.Latest,
		})
	}
	if ledgerState.History2.Latest == 0 {
		return ledgerState.Core.OldestOnStart, nil
	}

	return ledgerState.History2.Latest + 1, nil
}
