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

// GetMatchList - processes request to get the list of matches
func GetMatchList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetMatchList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getMatchListHandler{
		MatchQ: history2.NewMatchQ(historyRepo),
	}

	result, err := handler.GetMatchList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get match list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getMatchListHandler struct {
	MatchQ history2.MatchQ
	Log    *logan.Entry
}

// GetMatchList returns list of matches with related resources
func (h *getMatchListHandler) GetMatchList(request *requests.GetMatchList) (*regources.MatchListResponse, error) {
	q := h.MatchQ.Page(*request.PageParams).FilterByAssetPair(request.Filters.BaseAsset, request.Filters.QuoteAsset)

	coreMatches, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get matches from db")
	}

	response := regources.MatchListResponse{
		Data: make([]regources.Match, 0, len(coreMatches)),
	}

	for _, coreMatch := range coreMatches {
		response.Data = append(response.Data, resources.NewMatch(coreMatch))
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}
