// Package ledger provides useful utilities concerning ledgers within stellar,
// specifically a central location to store a cached snapshot of the state of
// both horizon's and stellar-core's view of the ledger.  This package is
// intended to be at the lowest levels of horizon's dependency tree, please keep
// it free of dependencies to other horizon packages.
package ledger

import (
	"sync"
	"time"

	"context"

	"github.com/saychao/horizon/log"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// SystemState represents a snapshot of both horizon's and stellar-core's view of the
// ledger.
type SystemState struct {
	Core     State `json:"core"`
	History  State `json:"history"`
	History2 State `json:"history_2"`
}

// State - represents a snapshot of currently available ledgers in the subsystem
type State struct {
	Latest                 int32     `json:"latest"`
	OldestOnStart          int32     `json:"oldest_on_start"`
	LastLedgerIncreaseTime time.Time `json:"last_ledger_increase_time"`
}

// CurrentState returns the cached snapshot of ledger state
// TODO: do not use static method
func CurrentState() SystemState {
	lock.RLock()
	defer lock.RUnlock()
	if instance == nil {
		panic(errors.New("trying to access current ledger state before StartLedgerStateUpdater"))
	}
	return *instance
}

type ledgerSeqProvider interface {
	// GetOldestLedgerSeq - returns oldest ledger sequence
	GetOldestLedgerSeq() (int32, error)
	// GetLatestLedgerSeq - returns latest ledger sequence available in DB
	GetLatestLedgerSeq() (int32, error)
}

// Config - helper struct which contains all resources needed by ledger state updater
type Config struct {
	CoreDB    string
	HistoryDB string
	Core      ledgerSeqProvider
	History   ledgerSeqProvider
	History2  ledgerSeqProvider
}

// StartLedgerStateUpdater - initializes current ledger state and start listeners for the updates in separate
// go routines
func StartLedgerStateUpdater(ctx context.Context, log *log.Entry, conf Config) error {
	err := initSystemState(conf)
	if err != nil {
		return errors.Wrap(err, "failed to init system state")
	}

	err = startNewListener(ctx, log.WithField("listener", "core_ledger_seq"), &lock, &instance.Core,
		"new_ledgers_seq", conf.CoreDB)
	if err != nil {
		return errors.Wrap(err, "failed to start core ledger seq listener")
	}

	err = startNewListener(ctx, log.WithField("listener", "history_ledger_seq"), &lock, &instance.History,
		"new_history_ledger_seq", conf.HistoryDB)
	if err != nil {
		return errors.Wrap(err, "failed to start history seq listener")
	}

	err = startNewListener(ctx, log.WithField("listener", "history2_ledger_seq"), &lock, &instance.History2,
		"new_ledgers_seq", conf.HistoryDB)
	if err != nil {
		return errors.Wrap(err, "failed to start history2 seq listener")
	}

	return nil
}

func initSystemState(conf Config) error {
	lock.Lock()
	defer lock.Unlock()
	instance = new(SystemState)
	var err error
	instance.Core, err = tryGetState(conf.Core)
	if err != nil {
		return errors.Wrap(err, "failed to get core state")
	}

	instance.History, err = tryGetState(conf.History)
	if err != nil {
		return errors.Wrap(err, "failed to get history state")
	}

	instance.History2, err = tryGetState(conf.History2)
	if err != nil {
		return errors.Wrap(err, "failed to get history2 state")
	}

	return nil
}

func tryGetState(provider ledgerSeqProvider) (State, error) {
	var result State
	var err error
	result.OldestOnStart, err = provider.GetOldestLedgerSeq()
	if err != nil {
		return State{}, errors.Wrap(err, "failed to get oldest ledger seq")
	}

	result.Latest, err = provider.GetLatestLedgerSeq()
	if err != nil {
		return State{}, errors.Wrap(err, "failed to get latest ledger seq")
	}

	result.LastLedgerIncreaseTime = time.Now().UTC()
	return result, nil
}

var instance *SystemState
var lock sync.RWMutex
