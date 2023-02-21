package operations

import "github.com/SafeRE-IT/horizon/db2/history2"

type updateDataHandler struct {
	effectsProvider
}

func (h *updateDataHandler) PermanentReject(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}

func (h *updateDataHandler) Fulfilled(details requestDetails) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(details.SourceAccountID)}, nil
}
