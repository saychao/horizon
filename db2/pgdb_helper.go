package db2

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func DeleteRange(
	db *pgdb.DB,
	start, end int64,
	table string,
	idCol string,
) error {
	del := sq.Delete(table).Where(
		fmt.Sprintf("%s >= ? AND %s < ?", idCol, idCol),
		start,
		end,
	)
	err := db.Exec(del)
	return err
}
