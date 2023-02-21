package storage

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

// Vote is helper struct to operate with `votes`
type Vote struct {
	repo *pgdb.DB
}

// NewVote - creates new instance of the `Vote`
func NewVote(repo *pgdb.DB) *Vote {
	return &Vote{
		repo: repo,
	}
}

// Insert - inserts new vote
func (q *Vote) Insert(vote history2.Vote) error {
	sql := sq.Insert("votes").
		Columns(
			"id", "voter_id", "poll_id", "data",
		).
		Values(
			vote.ID, vote.VoterID, vote.PollID, vote.VoteData,
		)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert vote", logan.F{
			"voter_id": vote.VoterID,
			"poll_id":  vote.PollID,
		})
	}

	return nil
}

// Update - updates existing vote
func (q *Vote) Update(vote history2.Vote) error {
	sql := sq.Update("votes").SetMap(map[string]interface{}{
		"data": vote.VoteData,
	}).Where("poll_id = ?", vote.PollID).Where("voter_id = ?", vote.VoterID)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update vote", logan.F{
			"voter_id": vote.VoterID,
			"poll_id":  vote.PollID,
		})
	}

	return nil
}

// Update - updates existing vote
func (q *Vote) Remove(voterID string, pollID uint64) error {
	sql := sq.Delete("votes").Where("voter_id = ?", voterID).Where("poll_id = ?", pollID)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to remove vote", logan.F{
			"voter_id": voterID,
			"poll_id":  pollID,
		})
	}

	return nil
}
