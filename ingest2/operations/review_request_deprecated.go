package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type deprecatedReviewRequestHandler struct {
}

// Fulfilled - always returns errors, as deprecated request must not occur in the core
func (h *deprecatedReviewRequestHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return nil, errors.From(errors.New("tried to ingest fulfilled deprecated reviewable request"), logan.F{
		"request_type": details.Request.Body.Type,
	})
}

// PermanentReject - always returns errors, as deprecated request must not occur in the core
func (h *deprecatedReviewRequestHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return nil, errors.From(errors.New("tried to ingest permanently rejected deprecated reviewable request"), logan.F{
		"request_type": details.Request.Body.Type,
	})
}
