package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// CreateData is helper struct to operate with `datas`
type Data struct {
	repo    *pgdb.DB
	updater sq.UpdateBuilder
}

// NewData - creates new instance of the `CreateData`
func NewData(repo *pgdb.DB) *Data {
	return &Data{
		repo:    repo,
		updater: sq.Update("data"),
	}
}

// Insert - inserts new data
func (q *Data) Insert(data history2.Data) error {

	sql := sq.Insert("data").
		Columns(
			"id", "type", "owner", "value",
		).
		Values(
			data.ID, data.Type, data.Owner, data.Value,
		)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert data", logan.F{"data_id": data.ID})
	}

	return nil
}

// Update - updates existing data
func (q *Data) Update(data history2.Data) error {
	sql := sq.Update("data").SetMap(map[string]interface{}{
		"owner": data.Owner,
		"type":  data.Type,
		"value": data.Value,
	}).Where(sq.Eq{"id": data.ID})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update data", logan.F{"data_id": data.ID})
	}

	return nil
}

// Remove - removes existing data
func (q *Data) Remove(dataID int64) error {
	sql := sq.Delete("data").Where(sq.Eq{"id": dataID})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to remove data", logan.F{"data_id": dataID})
	}

	return nil
}
