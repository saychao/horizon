package history2

import (
	"database/sql"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// PollsQ is a helper struct to aid in configuring queries that loads
// poll structures.
type PollsQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewPollsQ - creates new instance of PollsQ
func NewPollsQ(repo *pgdb.DB) PollsQ {
	return PollsQ{
		repo: repo,
		selector: sq.Select(
			"p.id",
			"p.permission_type",
			"p.number_of_choices",
			"p.data",
			"p.start_time",
			"p.end_time",
			"p.owner_id",
			"p.result_provider_id",
			"p.vote_confirmation_required",
			"p.details",
			"p.state",
		).From("polls p"),
	}
}

// GetByID loads a row from `polls`, by ID
// returns nil, nil - if poll does not exists
func (q PollsQ) GetByID(id int64) (*Poll, error) {
	q.selector = q.selector.Where("p.id = ?", id)
	return q.Get()
}

// FilterByOwner - returns q with filter by Owner
func (q PollsQ) FilterByOwner(ownerID string) PollsQ {
	q.selector = q.selector.Where("p.owner_id = ?", ownerID)
	return q
}

// FilterByResultProvider - returns q with filter by Owner
func (q PollsQ) FilterByResultProvider(resultProviderID string) PollsQ {
	q.selector = q.selector.Where("p.result_provider_id = ?", resultProviderID)
	return q
}

// FilterByPermissionType - returns q with filter by BaseAsset
func (q PollsQ) FilterByPermissionType(permissionType uint32) PollsQ {
	q.selector = q.selector.Where("p.permission_type = ?", permissionType)
	return q
}

// FilterByMaxEndTime - returns q with filter by max end time
func (q PollsQ) FilterByMaxEndTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("p.end_time <= ?", time)
	return q
}

// FilterByMaxStartTime - returns q with filter by start_time
func (q PollsQ) FilterByMaxStartTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("p.start_time <= ?", time)
	return q
}

// FilterByMinStartTime - returns q with filter by start_time
func (q PollsQ) FilterByMinStartTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("p.start_time >= ?", time)
	return q
}

// FilterByMinEndTime - returns q with filter by end_time
func (q PollsQ) FilterByMinEndTime(time time.Time) PollsQ {
	q.selector = q.selector.Where("p.end_time >= ?", time)
	return q
}

// FilterByPollType - returns q with filter by type
func (q PollsQ) FilterByPollType(pollType uint64) PollsQ {
	q.selector = q.selector.Where("data->'type'->>'value' = ?", pollType)
	return q
}

// FilterByPollState - returns q with filter by state
func (q PollsQ) FilterByState(state int32) PollsQ {
	q.selector = q.selector.Where("p.state = ?", state)
	return q
}

// FilterByVoteConfirmationRequired- returns q with filter by voteConfirmationRequired
func (q PollsQ) FilterByVoteConfirmationRequired(voteConfirmation bool) PollsQ {
	q.selector = q.selector.Where("p.vote_confirmation_required = ?", voteConfirmation)
	return q
}

// Page - returns Q with specified limit and offset params
func (q PollsQ) Page(params pgdb.CursorPageParams) PollsQ {
	q.selector = params.ApplyTo(q.selector, "p.id")
	return q
}

// Get - loads a row from `polls`
// returns nil, nil - if poll does not exists
// returns error if more than one CreatePoll found
func (q PollsQ) Get() (*Poll, error) {
	var result Poll
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load poll")
	}

	return &result, nil
}

//Select - selects slice from the db, if no polls found - returns nil, nil
func (q PollsQ) Select() ([]Poll, error) {
	var result []Poll
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load polls")
	}

	return result, nil
}
