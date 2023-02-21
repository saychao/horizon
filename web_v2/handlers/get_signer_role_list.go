package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetSignerRoleList - processes request to get the list of signerRoles
func GetSignerRoleList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getSignerRoleListHandler{
		SignerRolesQ: core2.NewSignerRoleQ(coreRepo),
		SignerRulesQ: core2.NewSignerRuleQ(coreRepo),
		Log:          ctx.Log(r),
	}

	request, err := requests.NewGetSignerRoleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSignerRoleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get signerRole list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSignerRoleListHandler struct {
	SignerRolesQ core2.SignerRoleQ
	SignerRulesQ core2.SignerRuleQ
	Log          *logan.Entry
}

// GetSignerRoleList returns the list of signerRoles with related resources
func (h *getSignerRoleListHandler) GetSignerRoleList(request *requests.GetSignerRoleList) (*regources.SignerRoleListResponse, error) {
	signerRoles, err := h.SignerRolesQ.Page(request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get signer role list")
	}

	response := &regources.SignerRoleListResponse{
		Data:  make([]regources.SignerRole, 0, len(signerRoles)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, signerRole := range signerRoles {
		signerRoleResponse := resources.NewSignerRole(signerRole)
		var rules []core2.SignerRule
		rules, err = h.SignerRulesQ.FilterByRole(signerRole.ID).Select()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load signer rules for role", logan.F{
				"role_id": signerRole.ID,
			})
		}

		for _, rule := range rules {
			ruleResponse := resources.NewSignerRule(rule)
			signerRoleResponse.Relationships.Rules.Data = append(signerRoleResponse.Relationships.Rules.Data,
				ruleResponse.Key)

			if request.ShouldInclude(requests.IncludeTypeRoleRules) {
				response.Included.Add(&ruleResponse)
			}
		}

		response.Data = append(response.Data, signerRoleResponse)
	}

	return response, nil
}
