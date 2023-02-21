package requests

import "net/http"

type GetParticipantEffect struct {
	ID uint64
	*base
}

func NewGetParticipantEffect(r *http.Request) (*GetParticipantEffect, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeHistoryAsset:            {},
			IncludeTypeHistoryOperationDetails: {},
			IncludeTypeHistoryEffect:           {},
			IncludeTypeHistoryOperation:        {},
		},
		supportedFilters: map[string]struct{}{},
	})
	if err != nil {
		return nil, err
	}

	request := GetParticipantEffect{base: b}

	request.ID, err = b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &request, nil
}
