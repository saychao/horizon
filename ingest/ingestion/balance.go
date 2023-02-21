package ingestion

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (ingest *Ingestion) TryIngestBalance(
	balanceID, asset, accountID string) (bool, error) {
	result, err := ingest.DB.ExecWithResult(squirrel.Expr(`
		insert into history_balances (balance_id, asset, account_id)
		values ($1, $2, $3) on conflict do nothing`,
		balanceID, asset, accountID))
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "failed to get rows affected")
	}
	return rows > 0, nil
}

func (ingest *Ingestion) TryIngestBalanceUpdate(
	balanceID string, amount, closeTime int64) error {
	err := ingest.DB.Exec(squirrel.
		Insert("history_balance_updates").
		SetMap(map[string]interface{}{
			"balance_id": balanceID,
			"amount":     amount,
			"updated_at": time.Unix(closeTime, 0).UTC(),
		}))
	return err
}
