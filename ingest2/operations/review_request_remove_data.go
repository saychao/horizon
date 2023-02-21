package operations

import "github.com/SafeRE-IT/horizon/db2/history2"

type removeDataHandler struct {
	effectsProvider
}

func (h *removeDataHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}

func (h *removeDataHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}
