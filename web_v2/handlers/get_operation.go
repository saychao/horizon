package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/web_v2/resources"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

// GetOperation - processes request to get operation and it's details by it's ID
func GetOperation(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getOperationHandler{
		OperationQ: history2.NewOperationQ(historyRepo),
		Log:        ctx.Log(r),
	}

	request, err := requests.NewGetOperation(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	result, err := handler.GetOperation(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get operation")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
	}

	ape.Render(w, result)
}

type getOperationHandler struct {
	OperationQ history2.OperationQ
	Log        *logan.Entry
}

func (h *getOperationHandler) GetOperation(request *requests.GetOperation) (*regources.OperationResponse, error) {
	op, err := h.OperationQ.FilterByID(request.ID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load operation")
	}

	if op == nil {
		return nil, nil
	}

	result := regources.OperationResponse{
		Data: resources.NewOperation(*op),
	}

	if request.ShouldInclude(requests.IncludeTypeOperationOperationDetails) {
		result.Included.Add(resources.NewOperationDetails(*op))
	}

	return &result, nil
}
