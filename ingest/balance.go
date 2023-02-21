package ingest

import (
	"github.com/pkg/errors"
	"github.com/saychao/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

func balanceUpdated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	balance := ledgerEntry.Data.Balance
	if balance == nil {
		return errors.New("expected non nil balance")
	}

	// seems like we have partial history, ensuring balance exists
	if !is.Paranoid {
		return nil
	}

	var b core.Balance
	err := is.Cursor.CoreQ().BalanceByID(&b, balance.BalanceId.AsString())
	if err != nil {
		return errors.Wrap(err, "failed to get balance")
	}
	_, err = is.Ingestion.TryIngestBalance(b.BalanceID, b.Asset, b.AccountID)
	if err != nil {
		return errors.Wrap(err, "failed to ingest balance")
	}

	return nil
}

func balanceCreated(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	balance := ledgerEntry.Data.Balance
	if balance == nil {
		return errors.New("Expected balance not to be nil")
	}

	_, err := is.Ingestion.TryIngestBalance(balance.BalanceId.AsString(),
		string(balance.Asset),
		balance.AccountId.Address())
	if err != nil {
		return errors.Wrap(err, "failed to ingest balance")
	}

	if err := balanceUpdated(is, ledgerEntry); err != nil {
		return errors.Wrap(err, "failed to updated balance")
	}
	return nil
}
