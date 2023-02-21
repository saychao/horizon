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

// GetAccountRule - processes request to get accountRule and it's details by accountRule code
func GetAccountRule(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountRuleHandler{
		AccountRulesQ: core2.NewAccountRuleQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetAccountRule(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountRule(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account rule", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getAccountRuleHandler struct {
	AccountRulesQ core2.AccountRuleQ
	Log           *logan.Entry
}

// GetAccountRule returns accountRule with related resources
func (h *getAccountRuleHandler) GetAccountRule(request *requests.GetAccountRule) (*regources.AccountRuleResponse, error) {
	accountRule, err := h.AccountRulesQ.FilterByID(request.ID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get account rule by id")
	}
	if accountRule == nil {
		return nil, nil
	}

	accountRuleResponse := resources.NewAccountRule(*accountRule)
	return &regources.AccountRuleResponse{
		Data: accountRuleResponse,
	}, nil
}
