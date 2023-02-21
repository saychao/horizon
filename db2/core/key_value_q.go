package core

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var selectKeyValue = sq.Select("kv.key", "kv.value").From("key_value_entry kv")

// KeyValueQI - provides methods to operate key-value
type KeyValueQI interface {
	// ByKey - selects KeyValue by key. Returns nil, nil if not found
	ByKey(key string) (*KeyValue, error)
	Select() ([]KeyValue, error)
}

type KeyValueQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *KeyValueQ) ByKey(key string) (*KeyValue, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := q.sql.Where("key = ?", key)

	var result KeyValue
	err := q.parent.Get(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load key value")
	}

	return &result, nil
}

// Select selects all existing KeyValues. Returns nil, nil if not found
func (q KeyValueQ) Select() ([]KeyValue, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []KeyValue
	err := q.parent.Select(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load all key values")
	}

	return result, nil
}
