package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetAccountRuleList - processes request to get the list of accountRules
func GetAccountRuleList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountRuleListHandler{
		AccountRulesQ: core2.NewAccountRuleQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetAccountRuleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountRuleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get accountRule list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAccountRuleListHandler struct {
	AccountRulesQ core2.AccountRuleQ
	Log           *logan.Entry
}

// GetAccountRuleList returns the list of accountRules with related resources
func (h *getAccountRuleListHandler) GetAccountRuleList(request *requests.GetAccountRuleList) (*regources.AccountRuleListResponse, error) {
	accountRules, err := h.AccountRulesQ.Page(request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get account rule list")
	}

	response := &regources.AccountRuleListResponse{
		Data:  make([]regources.AccountRule, 0, len(accountRules)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, accountRule := range accountRules {
		accountRuleResponse := resources.NewAccountRule(accountRule)
		response.Data = append(response.Data, accountRuleResponse)
	}

	return response, nil
}
