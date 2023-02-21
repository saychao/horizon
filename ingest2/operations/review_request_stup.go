package operations

import "github.com/saychao/horizon/db2/history2"

type reviewableRequestHandlerStub struct {
	effectsProvider
}

// Fulfilled - returns participant of fully approved request
func (h *reviewableRequestHandlerStub) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}

// PermanentReject - returns participants of fully rejected request
func (h *reviewableRequestHandlerStub) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}
