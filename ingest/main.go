// Package ingest contains the ingestion system for horizon.  This system takes
// data produced by the connected stellar-core database, transforms it and
// inserts it into the horizon database.
package ingest

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"sync"
	"time"

	"github.com/SafeRE-IT/horizon/corer"
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/ingest/ingestion"
	"github.com/SafeRE-IT/horizon/log"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	// CurrentVersion reflects the latest version of the ingestion
	// algorithm. As rows are ingested into the horizon database, this version is
	// used to tag them.  In the future, any breaking changes introduced by a
	// developer should be accompanied by an increase in this value.
	//
	// Scripts, that have yet to be ported to this codebase can then be leveraged
	// to re-ingest old data with the new algorithm, providing a seamless
	// transition when the ingested data's structure changes.
	CurrentVersion = 8
)

// Cursor iterates through a stellar core database's ledgers
type Cursor struct {
	// FirstLedger is the beginning of the range of ledgers (inclusive) that will
	// attempt to be ingested in this session.
	FirstLedger int32
	// LastLedger is the end of the range of ledgers (inclusive) that will
	// attempt to be ingested in this session.
	LastLedger int32
	// DB is the stellar-core db that data is ingested from.
	CoreDB *pgdb.DB

	Metrics *IngesterMetrics

	lg   int32
	tx   int
	op   int
	data *LedgerBundle
}

func (c *Cursor) CoreQ() core.QInterface {
	return &core.Q{DB: c.CoreDB}
}

func (c *Cursor) LedgerCloseTime() time.Time {
	return time.Unix(int64(c.data.Header.CloseTime), 0).UTC()
}

// LedgerBundle represents a single ledger's worth of novelty created by one
// ledger close
type LedgerBundle struct {
	Sequence        int32
	Header          core.LedgerHeader
	TransactionFees []core.TransactionFee
	Transactions    []core.Transaction
}

// System represents the data ingestion subsystem of horizon.
type System struct {
	// HorizonDB is the connection to the horizon database that ingested data will
	// be written to.
	HorizonDB *pgdb.DB

	// CoreDB is the stellar-core db that data is ingested from.
	CoreDB *pgdb.DB

	APIDB *pgdb.DB

	CoreConnector *corer.Connector

	Metrics IngesterMetrics

	CoreInfo *corer.Info

	lock    sync.Mutex
	current *Session
}

// IngesterMetrics tracks all the metrics for the ingestion subsystem
type IngesterMetrics struct {
	ClearLedgerTimer  metrics.Timer
	IngestLedgerTimer metrics.Timer
	LoadLedgerTimer   metrics.Timer
}

// Session represents a single attempt at ingesting data into the history
// database.
type Session struct {
	Cursor    *Cursor
	Ingestion *ingestion.Ingestion

	CoreInfo *corer.Info

	// ClearExisting causes the session to clear existing data from the horizon db
	// when the session is run.
	ClearExisting bool

	// Metrics is a reference to where the session should record its metric information
	Metrics *IngesterMetrics

	//
	// Results fields
	//
	// Ingested is the number of ledgers that were successfully ingested during
	// this session.
	Ingested int

	// Paranoid signals ingest routines to double-check everything and don't trust
	// history database state to be complete
	Paranoid bool

	CoreConnector *corer.Connector

	log *log.Entry
}

// New initializes the ingester, causing it to begin polling the stellar-core
// database for now ledgers and ingesting data into the horizon database.
func New(coreConnector *corer.Connector, coreInfo *corer.Info, core, horizon *pgdb.DB) *System {
	i := &System{
		CoreInfo:      coreInfo,
		HorizonDB:     horizon,
		CoreDB:        core,
		CoreConnector: coreConnector,
	}

	i.Metrics.ClearLedgerTimer = metrics.NewTimer()
	i.Metrics.IngestLedgerTimer = metrics.NewTimer()
	i.Metrics.LoadLedgerTimer = metrics.NewTimer()
	return i
}

// NewSession initialize a new ingestion session, from `first` to `last` using
// `i`.
func NewSession(paranoid bool, first, last int32, i *System) *Session {
	hdb := i.HorizonDB.Clone()
	coredb := i.CoreDB.Clone()

	return &Session{
		Paranoid: paranoid,
		Ingestion: &ingestion.Ingestion{
			DB:     hdb,
			CoreDB: coredb,
		},

		Cursor: &Cursor{
			FirstLedger: first,
			LastLedger:  last,
			CoreDB:      i.CoreDB,
			Metrics:     &i.Metrics,
		},
		CoreConnector: i.CoreConnector,

		CoreInfo: i.CoreInfo,
		Metrics:  &i.Metrics,
		log: log.WithField("service", "ingest_session").
			WithField("first_ledger", first).
			WithField("last_ledger", last),
	}
}
