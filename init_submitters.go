package horizon

import (
	"net/http"

	"time"

	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/log"
	hTxsub "github.com/saychao/horizon/txsub"
	"gitlab.com/distributed_lab/corer"
	"gitlab.com/distributed_lab/txsub"
)

func initSubmissionSystem(app *App) {
	logger := &log.WithField("service", "initSubmissionSystem").Entry
	cq := &core.Q{DB: app.CoreRepoLogged(logger)}
	hq := &history.Q{DB: app.HistoryRepoLogged(logger)}
	coreConnector, err := corer.NewConnector(&http.Client{
		Timeout: time.Duration(1 * time.Minute),
	}, app.config.StellarCoreURL)
	if err != nil {
		logger.WithError(err).Panic("Failed to create core connector")
	}
	app.submitter = &txsub.System{
		Pending:   txsub.NewDefaultSubmissionList(),
		Submitter: txsub.NewDefaultSubmitter(coreConnector),
		Results: &hTxsub.ResultsProvider{
			Core:    cq,
			History: hq,
		},
		NetworkPassphrase: app.CoreInfo.NetworkPassphrase,
	}
}

func init() {
	appInit.Add("txsub", initSubmissionSystem, "app-context", "log", "horizon-db", "core-db", "core-info")
}
