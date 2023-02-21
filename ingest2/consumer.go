package ingest2

import (
	"context"
	core "github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/log"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type corer interface {
	// SetCursor reports to core last ledger been processed, so that core could release resources during maintenance
	SetCursor(id string, lastLedger int32) error
}

type dbTxManager interface {
	Transaction(transactionFunc pgdb.TransactionFunc) error
}

// Handler - handles ledger and transactions applied for this ledger
type Handler interface {
	// Handle - processes ledger and stores corresponding data to db
	Handle(header *core.LedgerHeader, txs []core.Transaction) error
	// Name - returns name of the Handler
	Name() string
}

// Consumer - consumes ingest data and populate db with it
type Consumer struct {
	log         *logan.Entry
	dbTxManager dbTxManager
	dataSource  <-chan LedgerBundle
	handlers    []Handler
	corer       corer
}

// NewConsumer - creates new instance of consumer
func NewConsumer(log *logan.Entry, dbTxManager dbTxManager, corer corer, handlers []Handler,
	dataSource <-chan LedgerBundle) *Consumer {
	return &Consumer{
		log:         log.WithField("service", "ingest_consumer"),
		dbTxManager: dbTxManager,
		dataSource:  dataSource,
		handlers:    handlers,
		corer:       corer,
	}
}

// Start - starts consumer in separate goroutine. Must only be used once
func (c *Consumer) Start(ctx context.Context) {
	go c.run(ctx)
}

func (c *Consumer) run(ctx context.Context) {
	for {
		bundles := c.readBatch(ctx)
		if len(bundles) == 0 {
			c.log.Info("Fetched empty ledger bundle batch. It's clear sign that we are going to stop")
			return
		}

		localLog := c.log.WithFields(logan.F{
			"batch_len": len(bundles),
			"from":      bundles[0].Header.Sequence,
			"to":        bundles[len(bundles)-1].Header.Sequence,
		})

		localLog.Info("Starting to process new ledger bundles batch")
		running.UntilSuccess(ctx, localLog, "ingest_consumer", func(ctx context.Context) (bool, error) {
			err := c.processBatch(ctx, bundles)
			if err != nil {
				return false, err
			}

			return true, nil
		}, minErrorRecoveryPeriod, maxErrorRecoveryPeriod)
		localLog.Info("Ledger bundles batch processed")
	}

}

func (c *Consumer) readBatch(ctx context.Context) []LedgerBundle {
	const maxBatchSize = 500
	bundles := c.readAtLeastOne(ctx)
	for {
		select {
		case ledgerBundle, ok := <-c.dataSource:
			if !ok {
				return bundles
			}

			bundles = append(bundles, ledgerBundle)
			if len(bundles) >= maxBatchSize {
				return bundles
			}

		case <-ctx.Done():
			return nil
		default:
			return bundles
		}
	}
}

func (c *Consumer) readAtLeastOne(ctx context.Context) []LedgerBundle {
	select {
	case bundle, ok := <-c.dataSource:
		{
			if !ok {
				return nil
			}

			bundles := make([]LedgerBundle, 0, len(c.dataSource)+1)
			bundles = append(bundles, bundle)
			return bundles
		}
	case <-ctx.Done():
		return nil

	}
}

func (c *Consumer) processBatch(ctx context.Context, bundles []LedgerBundle) error {

	err := c.dbTxManager.Transaction(func() (err error) {
		for _, bundle := range bundles {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			err = c.processLedgerBundle(ctx, bundle)
			if err != nil {
				return errors.Wrap(err, "failed to process ledger bundle",
					log.F{"ledger_seq": bundle.Header.Sequence})
			}
		}
		return
	})

	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}

func (c *Consumer) processLedgerBundle(ctx context.Context, bundle LedgerBundle) error {
	for _, handler := range c.handlers {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		err := handler.Handle(&bundle.Header, bundle.Transactions)
		if err != nil {
			return errors.Wrap(err, "failed to handle ledger", logan.F{
				"ledger_seq":   bundle.Header.Sequence,
				"handler_name": handler.Name(),
			})
		}
	}

	return nil
}
