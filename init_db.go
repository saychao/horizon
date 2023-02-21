package horizon

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/log"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func initHorizonDb(app *App) {
	repo, err := pgdb.Open(pgdb.Opts{
		URL:                app.config.DatabaseURL,
		MaxOpenConnections: 12,
		MaxIdleConnections: 4,
	})

	if err != nil {
		log.Panic(err)
	}

	app.historyQ = &history.Q{DB: repo}
}

func initCoreDb(app *App) {
	repo, err := pgdb.Open(pgdb.Opts{
		URL:                app.config.StellarCoreDatabaseURL,
		MaxOpenConnections: 12,
		MaxIdleConnections: 4,
	})

	if err != nil {
		log.Panic(err)
	}

	app.coreQ = core.NewQ(repo)
}

func init() {
	appInit.Add("horizon-db", initHorizonDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
