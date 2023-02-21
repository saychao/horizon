package core

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var selectStatisticsV2 = sq.Select("st.id",
	"st.account_id",
	"st.stats_op_type",
	"st.asset_code",
	"st.is_convert_needed",
	"st.daily_out",
	"st.weekly_out",
	"st.monthly_out",
	"st.annual_out",
	"st.updated_at").From("statistics_v2 st")

type StatisticsV2QI interface {
	ForAccount(accountID string) ([]StatisticsV2Entry, error)
}

type StatisticsV2Q struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *StatisticsV2Q) ForAccount(accountID string) ([]StatisticsV2Entry, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := selectStatisticsV2.Where("st.account_id = ?", accountID)
	var result []StatisticsV2Entry
	err := q.parent.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load statistics_v2 for account")
	}

	for _, statisticsV2 := range result {
		statisticsV2.ClearObsolete(time.Now().UTC())
	}

	return result, nil
}
