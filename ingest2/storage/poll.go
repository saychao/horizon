package storage

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// CreatePoll is helper struct to operate with `polls`
type Poll struct {
	repo    *pgdb.DB
	updater sq.UpdateBuilder
}

// NewPoll - creates new instance of the `CreatePoll`
func NewPoll(repo *pgdb.DB) *Poll {
	return &Poll{
		repo:    repo,
		updater: sq.Update("polls"),
	}
}

// Insert - inserts new poll
func (q *Poll) Insert(poll history2.Poll) error {

	sql := sq.Insert("polls").
		Columns(
			"id", "permission_type", "number_of_choices", "data", "start_time",
			"end_time", "owner_id", "result_provider_id",
			"vote_confirmation_required", "details", "state",
		).
		Values(
			poll.ID, poll.PermissionType, poll.NumberOfChoices, poll.Data, poll.StartTime,
			poll.EndTime, poll.OwnerID, poll.ResultProviderID, poll.VoteConfirmationRequired,
			poll.Details, poll.State,
		)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert poll", logan.F{"poll_id": poll.ID})
	}

	return nil
}

// Update - updates existing poll
func (q *Poll) Update(poll history2.Poll) error {
	sql := sq.Update("polls").SetMap(map[string]interface{}{
		"permission_type":            poll.PermissionType,
		"number_of_choices":          poll.NumberOfChoices,
		"data":                       poll.Data,
		"start_time":                 poll.StartTime,
		"end_time":                   poll.EndTime,
		"owner_id":                   poll.OwnerID,
		"result_provider_id":         poll.ResultProviderID,
		"vote_confirmation_required": poll.VoteConfirmationRequired,
		"details":                    poll.Details,
		"state":                      poll.State,
	}).Where("id = ?", poll.ID)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{"poll_id": poll.ID})
	}

	return nil
}

func (q *Poll) SetState(state regources.PollState) *Poll {
	q.updater = q.updater.Set("state", state)
	return q
}

func (q *Poll) SetDetails(details json.RawMessage) *Poll {
	q.updater = q.updater.Set("details", details)
	return q
}

func (q *Poll) UpdateWhere(pollID uint64, shouldResetUpdater bool) error {
	q.updater = q.updater.Where(sq.Eq{"id": pollID})
	err := q.repo.Exec(q.updater)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{"poll_id": pollID})
	}

	if shouldResetUpdater {
		q.updater = sq.Update("polls")
	}

	return nil
}
