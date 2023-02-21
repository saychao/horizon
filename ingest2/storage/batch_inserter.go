package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/cheekybits/genny/generic"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//go:generate genny -in=$GOFILE -out=tx_batch_insert.go gen "valueType=history2.Transaction"
//go:generate genny -in=$GOFILE -out=participant_effect_batch_insert.go gen "valueType=history2.ParticipantEffect"
//go:generate genny -in=$GOFILE -out=ledger_changes_batch_insert.go gen "valueType=history2.LedgerChanges"
//go:generate genny -in=$GOFILE -out=operation_details_batch_insert.go gen "valueType=history2.Operation"

type valueType generic.Type

type valueTypeConvertToValues func(row valueType) []interface{}

func valueTypeBatchInsert(repo *pgdb.DB, rows []valueType, tableName string, columns []string, converter valueTypeConvertToValues) error {
	if len(rows) == 0 {
		return nil
	}

	sql := sq.Insert(tableName).Columns(columns...)

	paramsInQueue := 0
	for _, row := range rows {
		paramsInQueue += len(columns)
		if paramsInQueue > maxPostgresParams {
			err := repo.Exec(sql)
			if err != nil {
				return errors.Wrap(err, "failed to perform batch insert", logan.F{"rows_len": len(rows)})
			}

			sql = sq.Insert(tableName).Columns(columns...)
			paramsInQueue = len(columns)
		}

		sql = sql.Values(converter(row)...)
	}

	if paramsInQueue == 0 {
		return nil
	}

	err := repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to perform batch insert", logan.F{"rows_len": len(rows)})
	}

	return nil
}
