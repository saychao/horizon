package requests

import "net/http"

const (
	IncludeTypeOperationOperationDetails = "operation.details"
)

// GetOperation represents params to be specified by user for getOperation handler
type GetOperation struct {
	*base
	ID       uint64
	Includes struct {
		OperationDetails bool `include:"operation.details"`
	}
}

// NewGetOperation returns new instance of the GetOperation request
func NewGetOperation(r *http.Request) (*GetOperation, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeOperationOperationDetails: {},
		},
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("id")
	if err != nil {
		return nil, err
	}

	return &GetOperation{
		base: b,
		ID:   id,
	}, nil
}
