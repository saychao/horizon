package handlers

import (
	"net/http"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
)

// GetPollList - processes request to get the list of sales
func GetPollList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getPollListHandler{
		PollsQ: history2.NewPollsQ(historyRepo),
		Log:    ctx.Log(r),
	}

	request, err := requests.NewGetPollList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetPollList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get poll list ", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getPollListHandler struct {
	PollsQ history2.PollsQ
	Log    *logan.Entry
}

// GetPollList returns the list of assets with related resources
func (h *getPollListHandler) GetPollList(request *requests.GetPollList) (*regources.PollListResponse, error) {
	q := h.PollsQ

	if request.ShouldFilter(requests.FilterTypePollListOwner) {
		q = q.FilterByOwner(request.Filters.Owner)
	}

	if request.ShouldFilter(requests.FilterTypePollListResultProvider) {
		q = q.FilterByResultProvider(request.Filters.ResultProvider)
	}

	if request.ShouldFilter(requests.FilterTypePollListPermissionType) {
		q = q.FilterByPermissionType(request.Filters.PermissionType)
	}

	if request.ShouldFilter(requests.FilterTypePollListVoteConfirmation) {
		q = q.FilterByVoteConfirmationRequired(request.Filters.VoteConfirmation)
	}

	if request.ShouldFilter(requests.FilterTypePollListMaxEndTime) {
		q = q.FilterByMaxEndTime(*request.Filters.MaxEndTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListMaxStartTime) {
		q = q.FilterByMaxStartTime(*request.Filters.MaxStartTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListMinStartTime) {
		q = q.FilterByMinStartTime(*request.Filters.MinStartTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListMinEndTime) {
		q = q.FilterByMinEndTime(*request.Filters.MinEndTime)
	}

	if request.ShouldFilter(requests.FilterTypePollListState) {
		q = q.FilterByState(request.Filters.State)
	}

	if request.ShouldFilter(requests.FilterTypePollListPollType) {
		q = q.FilterByPollType(request.Filters.PollType)
	}

	historyPolls, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get poll list ")
	}

	response := &regources.PollListResponse{
		Data: make([]regources.Poll, 0, len(historyPolls)),
	}

	for _, historyPoll := range historyPolls {
		poll := resources.NewPoll(historyPoll)
		response.Data = append(response.Data, poll)
	}
	h.PopulateLinks(response, request)

	return response, nil
}

func (h *getPollListHandler) PopulateLinks(
	response *regources.PollListResponse, request *requests.GetPollList,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}
