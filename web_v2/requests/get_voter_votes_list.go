package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	IncludeTypeVotersVoteListAccount = "account"
	IncludeTypeVoterVoteListPolls    = "polls"
)

type GetVoterVoteList struct {
	*base
	VoterID    string
	PageParams *pgdb.CursorPageParams
	Includes   struct {
		Account bool `include:"account"`
		Polls   bool `include:"polls"`
	}
}

func NewGetVotersVotes(r *http.Request) (*GetVoterVoteList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeVoterVoteListPolls:    {},
			IncludeTypeVotersVoteListAccount: {},
		},
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	id := b.getString("voter")

	return &GetVoterVoteList{
		base:       b,
		PageParams: pageParams,
		VoterID:    id,
	}, nil
}
