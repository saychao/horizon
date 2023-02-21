package handlers

import (
	"net/http"

	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetMovements - processes request to get the list of participant effects
func GetMovements(w http.ResponseWriter, r *http.Request) {
	handler := newHistoryHandler(r)

	request, ok := handler.prepare(w, r)
	if !ok {
		return
	}

	result, err := handler.GetMovements(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get participant effect list", logan.F{})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

// GetMovements returns the list of participant effects with related resources
func (h *getHistory) GetMovements(request *requests.GetHistory) (regources.ParticipantsEffectListResponse, error) {
	q := h.ApplyFilters(request, h.EffectsQ).Movements()
	result, err := h.SelectAndPopulate(request, q, h.AssetsQ)
	if err != nil {
		return result, errors.Wrap(err, "failed to select and populate participant effects")
	}

	return result, nil
}
