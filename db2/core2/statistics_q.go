package core2

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type StatsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

func NewStatsQ(repo *pgdb.DB) StatsQ {
	return StatsQ{
		repo: repo,
		selector: sq.Select(
			"s.id",
			"s.account_id",
			"s.stats_op_type",
			"s.asset_code",
			"s.is_convert_needed",
			"s.daily_out",
			"s.weekly_out",
			"s.monthly_out",
			"s.annual_out",
			"s.updated_at").
			From("statistics_v2 s"),
	}
}

func (q *StatsQ) FilterByAccount(accountID string) *StatsQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": accountID})
	return q
}

func (q StatsQ) Select() ([]Statistics, error) {
	var stats []Statistics

	err := q.repo.Select(&stats, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get statistics from core db")
	}

	return stats, nil
}
