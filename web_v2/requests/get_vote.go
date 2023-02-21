package requests

import "net/http"

type GetVote struct {
	*base
	VoterID string
	PollID  int64
}

func NewGetVote(r *http.Request) (*GetVote, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	voter := b.getString("voter")
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetVote{
		base:    b,
		VoterID: voter,
		PollID:  int64(id),
	}, nil
}
