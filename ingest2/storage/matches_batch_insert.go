package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type history2MatchesToValues func(row history2.Match) []interface{}

func matchesBatchInsert(repo *pgdb.DB, rows []history2.Match, tableName string, columns []string, converterFn history2MatchesToValues) error {
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
			paramsInQueue = 0
		}

		sql = sql.Values(converterFn(row)...)
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
