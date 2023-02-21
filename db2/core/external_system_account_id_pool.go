package core

import "github.com/Masterminds/squirrel"

type ExternalSystemAccountIDPoolQ struct {
	parent *Q
}

func NewExternalSystemAccountIDPoolQ(parent *Q) *ExternalSystemAccountIDPoolQ {
	return &ExternalSystemAccountIDPoolQ{
		parent,
	}
}

func (q ExternalSystemAccountIDPoolQ) EntitiesCount() (map[string]uint64, error) {
	stmt := squirrel.Select("external_system_type", "count(1)").
		From("external_system_account_id_pool").
		Where("is_deleted = false").
		GroupBy("external_system_type")

	var dest []struct {
		Type  string `db:"external_system_type"`
		Count uint64 `db:"count"`
	}

	err := q.parent.Select(&dest, stmt)
	if err != nil {
		return nil, err
	}

	result := map[string]uint64{}
	for _, system := range dest {
		result[system.Type] = system.Count
	}

	return result, nil
}
