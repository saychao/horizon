package requests

import "net/http"

const (
	// IncludeTypePollParticipation - defines if participation should be included in the response
	IncludeTypePollParticipation = "participation"
	// IncludeTypePollParticipation - defines if votes should be included in the response
	IncludeTypePollParticipationVotes = "participation.votes"
)

var includeTypePollAll = map[string]struct{}{
	IncludeTypePollParticipation:      {},
	IncludeTypePollParticipationVotes: {},
}

type GetPoll struct {
	*base
	ID int64
}

func NewGetPoll(r *http.Request) (*GetPoll, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypePollAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}
	return &GetPoll{
		base: b,
		ID:   int64(id),
	}, nil
}
