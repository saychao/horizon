package requests

import "net/http"

const (
	// IncludeTypeReviewableRequestDetails - defines if we need to include request details to response
	IncludeTypeReviewableRequestDetails = "request_details"
)

var includeTypeReviewableRequestAll = map[string]struct{}{
	IncludeTypeReviewableRequestDetails: {},
}

// GetReviewableRequest - represents params to be specified by user for Get Request handler
type GetReviewableRequest struct {
	*base
	ID       uint64
	Includes struct {
		RequestDetails bool `include:"request_details"`
	}
}

//NewGetReviewableRequest - returns new instance of GetRequest request
func NewGetReviewableRequest(r *http.Request) (*GetReviewableRequest, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeReviewableRequestAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("id")
	if err != nil {
		return nil, err
	}

	return &GetReviewableRequest{
		base: b,
		ID:   id,
	}, nil
}
