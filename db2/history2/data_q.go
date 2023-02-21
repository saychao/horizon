package history2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var dataColumns = []string{
	"id",
	"type",
	"value",
	"owner",
}

// DataQ is a helper struct to aid in configuring queries that loads datas
type DataQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewDataQ- creates new instance of DataQ
func NewDataQ(repo *pgdb.DB) DataQ {
	return DataQ{
		repo: repo,
		selector: sq.Select(
			"data.id",
			"data.type",
			"data.value",
			"data.owner",
		).From("data data"),
	}
}

// GetByCode - get data by code
func (q DataQ) GetByID(id int64) (*Data, error) {
	q.selector = q.selector.Where(sq.Eq{"data.id": id})
	return q.Get()
}

//FilterByOwner - gets data by owner address, returns nil, nil if one does not exist
func (q DataQ) FilterByOwner(address string) DataQ {
	q.selector = q.selector.Where(sq.Eq{"data.owner": address})
	return q
}

//FilterByType - gets data by security type, returns nil, nil if one does not exist
func (q DataQ) FilterByType(dataType int64) DataQ {
	q.selector = q.selector.Where(sq.Eq{"data.type": dataType})
	return q
}

//Get - selects data from db, returns nil, nil if one does not exists
func (q DataQ) Get() (*Data, error) {
	var result Data
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load data")
	}

	return &result, nil
}

func (q DataQ) Page(params *pgdb.CursorPageParams) DataQ {
	q.selector = params.ApplyTo(q.selector, "data.id")
	return q
}

//Select - selects slice from the db, if no data found - returns nil, nil
func (q DataQ) Select() ([]Data, error) {
	var result []Data
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load data")
	}

	return result, nil
}
