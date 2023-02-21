package horizon

import (
	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2"
	"github.com/saychao/horizon/ingest2/changes"
	"github.com/saychao/horizon/ingest2/operations"
	"github.com/saychao/horizon/ingest2/storage"
	"github.com/saychao/horizon/ledger"
	"github.com/saychao/horizon/log"
)

func initIngester2(app *App) {
	if !app.config.Ingest {
		return
	}

	ctx := app.ctx
	logger := log.DefaultLogger.Entry.WithField("service", "ingest_v2")
	coreRepo := app.CoreRepoLogged(logger)

	txProvider := struct {
		*core2.LedgerHeaderQ
		core2.TransactionQ
	}{
		LedgerHeaderQ: core2.NewLedgerHeaderQ(coreRepo),
		TransactionQ:  core2.NewTransactionQ(coreRepo),
	}

	hRepo := app.HistoryRepoLogged(logger)
	const batchSize = 500
	ledgersChan := ingest2.NewProducer(txProvider, history2.NewLedgerQ(hRepo), logger, batchSize, ledger.CurrentState).
		Start(ctx)

	accountStorage := storage.NewAccount(hRepo, coreRepo)
	balanceStorage := storage.NewBalance(hRepo, coreRepo, accountStorage)
	swapStorage := storage.NewSwap(hRepo)

	ledgerChangesHandler := changes.NewHandler(
		accountStorage,
		balanceStorage,
		storage.NewReviewableRequest(hRepo),
		storage.NewSale(hRepo),
		storage.NewAssetPair(hRepo),
		storage.NewPoll(hRepo),
		storage.NewVote(hRepo),
		storage.NewAccountSpecificRules(hRepo),
		accountStorage,
		swapStorage,
		storage.NewAsset(hRepo),
		storage.NewData(hRepo),
		storage.NewDeferredPayment(hRepo),
	)

	idProvider := struct {
		*storage.Account
		*storage.Balance
	}{
		Account: accountStorage,
		Balance: balanceStorage,
	}

	reviewableRequestsStorage := struct {
		history2.ReviewableRequestsQ
	}{
		history2.NewReviewableRequestsQ(hRepo),
	}

	opHandler := operations.NewOperationsHandler(
		storage.NewOperationDetails(hRepo),
		storage.NewOpParticipants(hRepo),
		&idProvider,
		balanceStorage,
		swapStorage,
		storage.NewDeferredPayment(hRepo),
		reviewableRequestsStorage,
	)

	matchesHandler := ingest2.NewMatchesSaver(
		storage.NewMatches(hRepo),
	)

	consumer := ingest2.NewConsumer(logger.WithField("service", "ingest_data_consumer"), hRepo, app.CoreConnector, []ingest2.Handler{
		ingest2.NewLedgerHandler(storage.NewLedger(hRepo)),
		ingest2.NewTxSaver(storage.NewTx(hRepo)),
		ingest2.NewLedgerChangesHandler(storage.NewLedgerChange(hRepo)),
		ledgerChangesHandler,
		matchesHandler,
		opHandler,
	}, ledgersChan)

	consumer.Start(ctx)
}

func init() {
	appInit.Add("ingester2", initIngester2, "app-context", "log", "horizon-db", "core-db", "core_connector", "core-info", "ledger-state")
}
