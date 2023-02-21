package handlers

import (
	"net/http"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// GetHistory - processes request to get the list of participant effects
func GetHistory(w http.ResponseWriter, r *http.Request) {
	handler := newHistoryHandler(r)

	request, ok := handler.prepare(w, r)
	if !ok {
		return
	}

	result, err := handler.GetHistory(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get participant effect list", logan.F{})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

// GetHistory returns the list of participant effects with related resources
func (h *getHistory) GetHistory(request *requests.GetHistory) (regources.ParticipantsEffectListResponse, error) {
	q := h.ApplyFilters(request, h.EffectsQ)
	result, err := h.SelectAndPopulate(request, q, h.AssetsQ)
	if err != nil {
		return result, errors.Wrap(err, "failed to select and populate participant effects")
	}

	return result, nil
}
