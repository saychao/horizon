package db2

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// PageQueryV2 represents a portion of a Query struct concerned with paging
// through a large dataset. Is used for page-based(offset, limit) pagination
type PageQueryV2 struct {
	Limit  uint64
	Offset uint64
	Page   uint64
}

const (
	// Default limit is the default amount of records will be returned
	DefaultLimit = 15
	MaxLimit     = 50
)

// NewPageQueryV2 accepts the page derived for request
// URL params and calculates the proper offset for the
// sql-select
func NewPageQueryV2(page uint64, limit uint64) (result PageQueryV2, err error) {
	if err != nil {
		err = errors.Wrap(err, "Got invalid page")
		return
	}

	switch {
	case limit == 0:
		result.Limit = DefaultLimit
		break
	case limit < 0:
		err = ErrInvalidLimit
		return
	case limit > MaxLimit:
		result.Limit = MaxLimit
		break
	default:
		result.Limit = limit
	}

	result.Offset = page * limit
	result.Page = page

	return
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method provides the default case for paging: int64
// cursor-based paging by an id column.
func (p PageQueryV2) ApplyTo(sql sq.SelectBuilder) (sq.SelectBuilder, error) {
	return sql.Limit(p.Limit).Offset(p.Offset), nil
}
