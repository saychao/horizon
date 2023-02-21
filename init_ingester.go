package horizon

import (
	"github.com/SafeRE-IT/horizon/ingest"
	"github.com/SafeRE-IT/horizon/log"
)

func initIngester(app *App) {
	if !app.config.Ingest {
		return
	}

	logger := log.DefaultLogger.Entry.WithField("service", "ingest")

	ingester := ingest.New(app.CoreConnector,
		app.CoreInfo,
		app.CoreRepoLogged(logger),
		app.HistoryRepoLogged(logger),
	)

	if err := ingester.IntegrityCheck(); err != nil {
		log.WithField("service", "ingester").
			WithError(err).
			Fatal("integrity check failed")
	}

	app.ingester = ingester
}

func init() {
	appInit.Add("ingester", initIngester, "core_connector", "app-context", "log", "horizon-db", "core-db", "core-info", "ledger-state")
}
