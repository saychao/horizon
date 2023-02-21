package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/history2"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetReviewableRequest - processes request to get reviewable request and it's details by it's ID
func GetReviewableRequest(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getReviewableRequestHandler{
		Log:       ctx.Log(r),
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
	}

	request, err := requests.NewGetReviewableRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRecord, err := handler.GetReviewableRequest(*request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get reviewable request", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if historyRecord == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, &historyRecord.Requestor, &historyRecord.Reviewer) {
		return
	}

	handler.RenderResponse(r, w, request, historyRecord)
}

type getReviewableRequestHandler struct {
	Log       *logan.Entry
	RequestsQ history2.ReviewableRequestsQ
}

// GetReviewableRequest returns gets request from history database
func (h *getReviewableRequestHandler) GetReviewableRequest(request requests.GetReviewableRequest) (*history2.ReviewableRequest, error) {
	historyRequest, err := h.RequestsQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get request by ID")
	}

	return historyRequest, nil
}

// RenderResponse renders reviewable request with related resources
func (h *getReviewableRequestHandler) RenderResponse(r *http.Request, w http.ResponseWriter,
	request *requests.GetReviewableRequest, historyRecord *history2.ReviewableRequest) {

	response := regources.ReviewableRequestResponse{
		Data: resources.NewRequest(*historyRecord),
	}

	if request.ShouldInclude(requests.IncludeTypeReviewableRequestDetails) {
		response.Included.Add(resources.NewRequestDetails(*historyRecord))
	}

	ape.Render(w, response)
}
