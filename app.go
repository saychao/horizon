package horizon

import (
	"fmt"
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
	"runtime"
	"sync"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/saychao/horizon/cache"
	"github.com/saychao/horizon/config"
	"github.com/saychao/horizon/corer"
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/ingest"
	"github.com/saychao/horizon/ledger"
	"github.com/saychao/horizon/log"
	"github.com/saychao/horizon/render/sse"
	txsub2 "github.com/saychao/horizon/txsub/v2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/txsub"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	graceful "gopkg.in/tylerb/graceful.v1"
)

// You can override this variable using: gb build -ldflags "-X main.version aabbccdd"
var version = ""

// App represents the root of the state of a horizon instance.
type App struct {
	config         config.Config
	web            *Web
	webV2          *WebV2
	historyQ       history.QInterface
	coreQ          core.QInterface
	ctx            context.Context
	cancel         func()
	submitter      *txsub.System
	submitterV2    *txsub2.System
	ingester       *ingest.System
	ticks          *time.Ticker
	CoreInfo       *corer.Info
	horizonVersion string
	cacheProvider  *cache.Provider
	CoreConnector  *corer.Connector

	// metrics
	metrics                  metrics.Registry
	historyLatestLedgerGauge metrics.Gauge
	historyElderLedgerGauge  metrics.Gauge
	horizonConnGauge         metrics.Gauge
	coreLatestLedgerGauge    metrics.Gauge
	coreElderLedgerGauge     metrics.Gauge
	coreConnGauge            metrics.Gauge
	goroutineGauge           metrics.Gauge
}

// SetVersion records the provided version string in the package level `version`
// var, which will be used for the reported horizon version.
func SetVersion(v string) {
	version = v
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.Config) (*App, error) {
	result := &App{config: config}
	result.horizonVersion = version
	result.ticks = time.NewTicker(1 * time.Second)
	result.init()
	return result, nil
}

// Serve starts the horizon web server, binding it to a socket, setting up
// the shutdown signals.
func (a *App) Serve() {
	a.web.router.Compile()
	http.Handle("/v3/", a.webV2.mux)
	http.Handle("/", a.web.router)

	addr := fmt.Sprintf(":%d", a.config.Port)

	srv := &graceful.Server{
		Timeout: 10 * time.Second,

		Server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},

		ShutdownInitiated: func() {
			log.Info("received signal, gracefully stopping")
			a.Close()
		},
	}

	http2.ConfigureServer(srv.Server, nil)

	log.Infof("Starting horizon on %s", addr)

	go a.run()

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	log.Info("stopped")
}

// Close cancels the app and forces the closure of db connections
func (a *App) Close() {
	a.cancel()
	a.ticks.Stop()

	a.historyQ.GetRepo().RawDB().Close()
	a.coreQ.GetRepo().RawDB().Close()

}

// CoreQ returns a helper object for performing sql queries against the
// stellar core database.
func (a *App) CoreQ() core.QInterface {
	return a.coreQ
}

// HistoryQ returns a helper object for performing sql queries against the
// history portion of horizon's database.
func (a *App) HistoryQ() history.QInterface {
	return a.historyQ
}

// CoreRepoLogged returns a new repo that loads data from the core database.
func (a *App) CoreRepoLogged(log *logan.Entry) *pgdb.DB {
	return a.coreQ.GetRepo().Clone()
}

// HistoryRepoLogged returns a new repo that loads data from the horizon database.
func (a *App) HistoryRepoLogged(log *logan.Entry) *pgdb.DB {
	return a.historyQ.GetRepo().Clone()
}

// IsHistoryStale returns true if the latest history ledger is more than
// `StaleThreshold` ledgers behind the latest core ledger
func (a *App) IsHistoryStale() bool {
	if a.config.StaleThreshold == 0 {
		return false
	}

	ls := ledger.CurrentState()
	return (ls.Core.Latest - ls.History.Latest) > int32(a.config.StaleThreshold)
}

// UpdateCoreInfo updates the value of coreVersion and networkPassphrase
// from the Stellar core API.
func (a *App) UpdateCoreInfo() error {
	if a.config.StellarCoreURL == "" {
		return nil
	}

	var info *corer.Info
	info, err := a.CoreConnector.GetCoreInfo()
	if err != nil {
		log.WithField("service", "core-info").WithError(err).Error("could not load stellar-core info")
		return errors.Wrap(err, "could not load stellar-core info")
	}

	a.CoreInfo = info

	return nil
}

// UpdateMetrics triggers a refresh of several metrics gauges, such as open
// db connections and ledger state
func (a *App) UpdateMetrics() {
	a.goroutineGauge.Update(int64(runtime.NumGoroutine()))
	ls := ledger.CurrentState()
	a.historyLatestLedgerGauge.Update(int64(ls.History.Latest))
	a.historyElderLedgerGauge.Update(int64(ls.History.OldestOnStart))
	a.coreLatestLedgerGauge.Update(int64(ls.Core.Latest))
	a.coreElderLedgerGauge.Update(int64(ls.Core.OldestOnStart))

	//a.horizonConnGauge.Update(int64(a.historyQ.Repo.DB.Stats().OpenConnections))
	//a.coreConnGauge.Update(int64(a.coreQ.Repo.DB.Stats().OpenConnections))
}

// UpdateWebV2Metrics updates the metrics for the web_v2 requests
func (a *App) UpdateWebV2Metrics(requestDuration time.Duration, responseStatus int) {
	a.webV2.metrics.Update(requestDuration, responseStatus)
}

// Tick triggers horizon to update all of it's background processes such as
// transaction submission, metrics, ingestion and reaping.
func (a *App) Tick() {
	var wg sync.WaitGroup
	log.Debug("ticking app")
	// update ledger state and stellar-core info in parallel
	wg.Add(1)
	go func() { a.UpdateCoreInfo(); wg.Done() }()
	wg.Wait()

	if a.ingester != nil {
		go a.ingester.Tick()
	}

	wg.Add(1)
	go func() { a.submitter.Tick(a.ctx); wg.Done() }()
	wg.Wait()

	sse.Tick()

	// finally, update metrics
	a.UpdateMetrics()
	log.Debug("finished ticking app")
}

// Init initializes app, using the config to populate db connections and
// whatnot.
func (a *App) init() {
	appInit.Run(a)
}

// run is the function that runs in the background that triggers Tick each
// second
func (a *App) run() {
	for {
		select {
		case <-a.ticks.C:
			a.Tick()
		case <-a.ctx.Done():
			log.Info("finished background ticker")
			return
		}
	}
}

func (a *App) Conf() config.Config {
	return a.config
}
