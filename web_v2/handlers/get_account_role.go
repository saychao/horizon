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

// GetAccountRole - processes request to get accountRole and it's details by accountRole code
func GetAccountRole(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountRoleHandler{
		AccountRolesQ: core2.NewAccountRoleQ(coreRepo),
		AccountRulesQ: core2.NewAccountRuleQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetAccountRole(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountRole(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account role", logan.F{
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

type getAccountRoleHandler struct {
	AccountRolesQ core2.AccountRoleQ
	AccountRulesQ core2.AccountRuleQ
	Log           *logan.Entry
}

// GetAccountRole returns accountRole with related resources
func (h *getAccountRoleHandler) GetAccountRole(request *requests.GetAccountRole) (*regources.AccountRoleResponse, error) {
	accountRole, err := h.AccountRolesQ.FilterByID(request.ID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get accountRole by code")
	}
	if accountRole == nil {
		return nil, nil
	}

	rules, err := h.AccountRulesQ.FilterByRole(request.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load account rules for role")
	}

	accountRoleResponse := resources.NewAccountRole(*accountRole)
	var included regources.Included
	for _, rule := range rules {
		ruleResponse := resources.NewAccountRule(rule)
		accountRoleResponse.Relationships.Rules.Data = append(accountRoleResponse.Relationships.Rules.Data,
			ruleResponse.Key)

		if request.ShouldInclude(requests.IncludeTypeRoleRules) {
			included.Add(&ruleResponse)
		}
	}

	return &regources.AccountRoleResponse{
		Data:     accountRoleResponse,
		Included: included,
	}, nil
}
