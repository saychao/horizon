package handlers

import (
	"net/http"
	"time"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/generator"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetPoll(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetPoll(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getPollHandler{
		LedgerHeaderQ: *core2.NewLedgerHeaderQ(ctx.CoreRepo(r)),
		VotesQ:        history2.NewVotesQ(ctx.HistoryRepo(r)),
		PollsQ:        history2.NewPollsQ(ctx.HistoryRepo(r)),
		Log:           ctx.Log(r),
	}

	result, err := handler.getPoll(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get poll", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if request.ShouldIncludeAny(
		requests.IncludeTypePollParticipation,
		requests.IncludeTypePollParticipationVotes,
	) {
		if !isAllowed(r, w, &result.Data.Relationships.Owner.Data.ID, &result.Data.Relationships.ResultProvider.Data.ID) {
			return
		}
	}

	ape.Render(w, result)
}

type getPollHandler struct {
	LedgerHeaderQ core2.LedgerHeaderQ
	PollsQ        history2.PollsQ
	VotesQ        history2.VotesQ
	Log           *logan.Entry
}

// GetSale returns sale with related resources
func (h *getPollHandler) getPoll(request *requests.GetPoll) (*regources.PollResponse, error) {

	record, err := h.PollsQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get poll")
	}

	if record == nil {
		return nil, nil
	}

	resource := resources.NewPoll(*record)
	response := &regources.PollResponse{
		Data: resource,
	}

	if !request.ShouldIncludeAny(
		requests.IncludeTypePollParticipation,
		requests.IncludeTypePollParticipationVotes,
	) {
		return response, nil
	}

	votes, err := h.VotesQ.FilterByPollID(request.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get votes for poll")
	}

	if request.ShouldInclude(requests.IncludeTypePollParticipation) {
		outcomeKey := resources.NewParticipation(record.ID, votes)
		response.Included.Add(&outcomeKey)
	}

	if request.ShouldInclude(requests.IncludeTypePollParticipationVotes) {
		ledgerHeaders, err := h.GetLedgerHeaders(votes)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get ledger headers")
		}

		for _, v := range votes {
			ledgerSequence := generator.GetSeqFromInt64(v.ID) // ledger sequence
			created := time.Unix(ledgerHeaders[ledgerSequence].CloseTime, 0)
			v.VoteData.CreationTime = &created
			vote := resources.NewVote(v)

			response.Included.Add(&vote)
		}
	}

	return response, nil
}

func (h *getPollHandler) GetLedgerHeaders(votes []history2.Vote) (map[int32]core2.LedgerHeader, error) {
	ledgerSeq := make([]int32, len(votes))
	for i, vote := range votes {
		ledgerSeq[i] = generator.GetSeqFromInt64(vote.ID)
	}

	ledgerHeaders, err := h.LedgerHeaderQ.SelectBySequence(ledgerSeq)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get header of ledger sequence")
	}

	result := make(map[int32]core2.LedgerHeader)
	for _, ledger := range ledgerHeaders {
		result[ledger.Sequence] = ledger
	}

	return result, nil
}
