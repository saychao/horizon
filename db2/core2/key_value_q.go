package core2

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type KeyValueQ struct {
	repo     *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewKeyValueQ(repo *pgdb.DB) *KeyValueQ {
	return &KeyValueQ{
		repo:     repo,
		selector: squirrel.Select("key", "value").From("key_value_entry"),
	}
}

func (q *KeyValueQ) ByKey(key string) (*KeyValue, error) {
	var result KeyValue
	stmt := q.selector.Where("key = ?", key)
	err := q.repo.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func (q *KeyValueQ) Page(params *pgdb.OffsetPageParams) *KeyValueQ {
	q.selector = params.ApplyTo(q.selector, "key")
	return q
}

func (q *KeyValueQ) Select() ([]KeyValue, error) {
	var result []KeyValue
	err := q.repo.Select(&result, q.selector)
	return result, err
}
