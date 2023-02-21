package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/generator"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

func GetVoterVotesList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)

	handler := getVoteListHandler{
		VotesQ:        history2.NewVotesQ(historyRepo),
		PollsQ:        history2.NewPollsQ(historyRepo),
		AccountsQ:     core2.NewAccountsQ(coreRepo),
		LedgerHeaderQ: *core2.NewLedgerHeaderQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetVotersVotes(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, &request.VoterID) {
		return
	}

	result, err := handler.GetVoterVotesList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get votes list")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

func (h *getVoteListHandler) GetVoterVotesList(request *requests.GetVoterVoteList) (*regources.VoteListResponse, error) {
	q := h.VotesQ.
		FilterByVoterID(request.VoterID).
		Page(*request.PageParams)

	historyVotes, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "cannot select votes")
	}

	response := regources.VoteListResponse{
		Data: make([]regources.Vote, 0, len(historyVotes)),
	}

	ledgerHeaders, err := h.GetLedgerHeaders(historyVotes)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get ledger headers")
	}

	for _, vote := range historyVotes {
		ledgerSequence := generator.GetSeqFromInt64(vote.ID) // ledger sequence
		created := time.Unix(ledgerHeaders[ledgerSequence].CloseTime, 0)
		vote.VoteData.CreationTime = &created

		response.Data = append(response.Data, resources.NewVote(vote))
		if request.ShouldInclude(requests.IncludeTypeVoterVoteListPolls) {
			historyPoll, err := h.PollsQ.GetByID(vote.PollID)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get voter account")
			}

			poll := resources.NewPoll(*historyPoll)
			response.Included.Add(&poll)
		}
	}

	if len(response.Data) > 0 {
		if request.ShouldInclude(requests.IncludeTypeVotersVoteListAccount) {
			coreAccount, err := h.AccountsQ.GetByAddress(request.VoterID)
			if err != nil {
				return nil, errors.Wrap(err, "cannot get voter account")
			}

			account := resources.NewAccount(*coreAccount, nil)
			response.Included.Add(&account)
		}

		response.Links = request.GetCursorLinks(*request.PageParams, fmt.Sprintf("%d", historyVotes[len(historyVotes)-1].ID))
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}

func (h *getVoteListHandler) GetLedgerHeaders(votes []history2.Vote) (map[int32]core2.LedgerHeader, error) {
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
