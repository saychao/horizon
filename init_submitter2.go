package horizon

import (
	"time"

	"github.com/SafeRE-IT/horizon/ingest2/storage"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/lib/pq"

	"github.com/SafeRE-IT/horizon/log"
	"github.com/SafeRE-IT/horizon/txsub/v2"
)

func initSubmissionV2(app *App) {
	logger := log.WithField("service", "submitter_v2")
	cq := core2.NewTransactionQ(app.CoreRepoLogged(&logger.Entry))
	hq := history2.NewTransactionsQ(app.HistoryRepoLogged(&logger.Entry))
	coreConnector := app.CoreConnector

	coreListener := pq.NewListener(
		app.config.StellarCoreDatabaseURL,
		time.Second,
		5*time.Second,
		log.PQEvent(logger),
	)

	histListener := pq.NewListener(
		app.config.DatabaseURL,
		time.Second,
		5*time.Second,
		log.PQEvent(logger),
	)
	err := coreListener.Listen(storage.ChanNewLedgerSeq)
	if err != nil {
		panic(errors.Wrap(err, "failed to init history listener", logan.F{
			"channel": storage.ChanNewLedgerSeq,
		}))
	}

	err = histListener.Listen(storage.ChanNewLedgerSeq)
	if err != nil {
		panic(errors.Wrap(err, "failed to init history listener", logan.F{
			"channel": storage.ChanNewLedgerSeq,
		}))
	}
	app.submitterV2 = &txsub.System{
		Log:               logger,
		SubmissionTimeout: time.Minute + 30*time.Second,
		List:              txsub.NewDefaultSubmissionList(10 * time.Second),
		Submitter:         txsub.NewDefaultSubmitter(*coreConnector),
		CoreListener:      coreListener,
		HistoryListener:   histListener,
		Results: &txsub.ResultsProvider{
			Core:    cq,
			History: hq,
		},
		NetworkPassphrase: app.CoreInfo.NetworkPassphrase,
	}

	app.submitterV2.Start(app.ctx)
}

func init() {
	appInit.Add("submitter_v2", initSubmissionV2, "app-context", "log", "horizon-db", "core-db", "core-info")
}
