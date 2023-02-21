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

// GetAccountRoleList - processes request to get the list of accountRoles
func GetAccountRoleList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountRoleListHandler{
		AccountRolesQ: core2.NewAccountRoleQ(coreRepo),
		AccountRulesQ: core2.NewAccountRuleQ(coreRepo),
		Log:           ctx.Log(r),
	}

	request, err := requests.NewGetAccountRoleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountRoleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get accountRole list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getAccountRoleListHandler struct {
	AccountRolesQ core2.AccountRoleQ
	AccountRulesQ core2.AccountRuleQ
	Log           *logan.Entry
}

// GetAccountRoleList returns the list of accountRoles with related resources
func (h *getAccountRoleListHandler) GetAccountRoleList(request *requests.GetAccountRoleList) (*regources.AccountRoleListResponse, error) {
	accountRoles, err := h.AccountRolesQ.Page(request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get account role list")
	}

	response := &regources.AccountRoleListResponse{
		Data:  make([]regources.AccountRole, 0, len(accountRoles)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, accountRole := range accountRoles {
		accountRoleResponse := resources.NewAccountRole(accountRole)
		var rules []core2.AccountRule
		rules, err = h.AccountRulesQ.FilterByRole(accountRole.ID).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load account rules for role", logan.F{
				"role_id": accountRole.ID,
			})
		}

		for _, rule := range rules {
			ruleResponse := resources.NewAccountRule(rule)
			accountRoleResponse.Relationships.Rules.Data = append(accountRoleResponse.Relationships.Rules.Data,
				ruleResponse.Key)

			if request.ShouldInclude(requests.IncludeTypeRoleRules) {
				response.Included.Add(&ruleResponse)
			}
		}

		response.Data = append(response.Data, accountRoleResponse)
	}

	return response, nil
}
