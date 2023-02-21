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

// GetSignerRole - processes request to get signerRole and it's details by signerRole code
func GetSignerRole(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getSignerRoleHandler{
		SignerRolesQ: core2.NewSignerRoleQ(coreRepo),
		SignerRulesQ: core2.NewSignerRuleQ(coreRepo),
		Log:          ctx.Log(r),
	}

	request, err := requests.NewGetSignerRole(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSignerRole(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get signer role", logan.F{
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

type getSignerRoleHandler struct {
	SignerRolesQ core2.SignerRoleQ
	SignerRulesQ core2.SignerRuleQ
	Log          *logan.Entry
}

// GetSignerRole returns signerRole with related resources
func (h *getSignerRoleHandler) GetSignerRole(request *requests.GetSignerRole) (*regources.SignerRoleResponse, error) {
	signerRole, err := h.SignerRolesQ.FilterByID(request.ID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get signerRole by code")
	}
	if signerRole == nil {
		return nil, nil
	}

	rules, err := h.SignerRulesQ.FilterByRole(request.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load signer rules for role")
	}

	signerRoleResponse := resources.NewSignerRole(*signerRole)
	var included regources.Included
	for _, rule := range rules {
		ruleResponse := resources.NewSignerRule(rule)
		signerRoleResponse.Relationships.Rules.Data = append(signerRoleResponse.Relationships.Rules.Data,
			ruleResponse.Key)

		if request.ShouldInclude(requests.IncludeTypeRoleRules) {
			included.Add(&ruleResponse)
		}
	}

	return &regources.SignerRoleResponse{
		Data:     signerRoleResponse,
		Included: included,
	}, nil
}
