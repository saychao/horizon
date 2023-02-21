package history2

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// VotesQ is a helper struct to aid in configuring queries that loads
// poll structures.
type VotesQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewVotesQ - creates new instance of VotesQ
func NewVotesQ(repo *pgdb.DB) VotesQ {
	return VotesQ{
		repo: repo,
		selector: sq.Select(
			"v.id",
			"v.poll_id",
			"v.voter_id",
			"v.data",
		).From("votes v"),
	}
}

// FilterByVoterID - returns q with filter by voter ID
func (q VotesQ) FilterByVoterID(voterID string) VotesQ {
	q.selector = q.selector.Where("v.voter_id = ?", voterID)
	return q
}

// FilterByPollID - returns q with filter by PollID
func (q VotesQ) FilterByPollID(pollID int64) VotesQ {
	q.selector = q.selector.Where("v.poll_id = ?", pollID)
	return q
}

// Page - returns Q with specified limit and offset params
func (q VotesQ) Page(params pgdb.CursorPageParams) VotesQ {
	q.selector = params.ApplyTo(q.selector, "v.id")
	return q
}

// Get - loads a row from `sales`
// returns nil, nil - if sale does not exists
// returns error if more than one Vote found
func (q VotesQ) Get() (*Vote, error) {
	var result Vote
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load vote")
	}

	return &result, nil
}

//Select - selects slice from the db, if no sales found - returns nil, nil
func (q VotesQ) Select() ([]Vote, error) {
	var result []Vote
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load votes")
	}

	return result, nil
}
