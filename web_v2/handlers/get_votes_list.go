package handlers

import (
	"fmt"
	"net/http"

	"github.com/saychao/horizon/db2/core2"

	"github.com/saychao/horizon/db2/history2"

	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetVoteList - processes request to get the list of sales
func GetVoteList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	handler := getVoteListHandler{
		VotesQ: history2.NewVotesQ(historyRepo),
		PollsQ: history2.NewPollsQ(historyRepo),
		Log:    ctx.Log(r),
	}

	request, err := requests.NewGetVoteList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	poll, err := handler.PollsQ.GetByID(request.PollID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get poll for vote", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if poll == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, &poll.OwnerID, &poll.ResultProviderID) {
		return
	}

	result, err := handler.GetVoteList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get vote list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getVoteListHandler struct {
	LedgerHeaderQ core2.LedgerHeaderQ
	VotesQ        history2.VotesQ
	PollsQ        history2.PollsQ
	AccountsQ     core2.AccountsQ
	Log           *logan.Entry
}

// GetVoteList returns the list of assets with related resources
func (h *getVoteListHandler) GetVoteList(request *requests.GetVoteList) (*regources.VoteListResponse, error) {
	q := h.VotesQ.FilterByPollID(request.PollID).Page(*request.PageParams)

	historyVotes, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get vote list")
	}

	response := &regources.VoteListResponse{
		Data: make([]regources.Vote, 0, len(historyVotes)),
	}

	for _, historyVote := range historyVotes {
		vote := resources.NewVote(historyVote)

		response.Data = append(response.Data, vote)
	}
	if len(historyVotes) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, fmt.Sprintf("%d", historyVotes[len(historyVotes)-1].ID))
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return response, nil
}
