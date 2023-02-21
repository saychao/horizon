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

// GetAccountSigners - processes request to get account signers
func GetAccountSigners(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountSignersHandler{
		SignersQ:    core2.NewSignerQ(coreRepo),
		AccountQ:    core2.NewAccountsQ(coreRepo),
		SignerRoleQ: core2.NewSignerRoleQ(coreRepo),
		SignerRuleQ: core2.NewSignerRuleQ(coreRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetAccountSigners(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetAccountSigners(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get account", logan.F{
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

type getAccountSignersHandler struct {
	SignersQ    core2.SignerQ
	AccountQ    core2.AccountsQ
	SignerRoleQ core2.SignerRoleQ
	SignerRuleQ core2.SignerRuleQ
	Log         *logan.Entry
}

// GetAccountSigners - returns signers for account
func (h *getAccountSignersHandler) GetAccountSigners(request *requests.GetAccountSigners) (*regources.SignerListResponse, error) {
	account, err := h.AccountQ.GetByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load account", logan.F{
			"account_address": request.Address,
		})
	}

	if account == nil {
		return nil, nil
	}

	signers, err := h.SignersQ.Page(request.PageParams).FilterByAccountAddress(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load signers for account", logan.F{
			"account_address": request.Address,
		})
	}

	response := regources.SignerListResponse{
		Data:  make([]regources.Signer, 0, len(signers)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, signerRaw := range signers {
		signer := resources.NewSigner(signerRaw)
		signer.Relationships.Role, err = h.getRole(request, &response.Included, signerRaw)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get role")
		}

		response.Data = append(response.Data, signer)
	}

	return &response, nil
}

func (h *getAccountSignersHandler) getRole(request *requests.GetAccountSigners,
	includes *regources.Included, signer core2.Signer,
) (*regources.Relation, error) {
	if !request.ShouldInclude(requests.IncludeTypeSignerRole) {
		role := resources.NewSignerRoleKey(signer.RoleID)
		return role.AsRelation(), nil
	}

	roleRaw, err := h.SignerRoleQ.GetByID(signer.RoleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer role")
	}

	if roleRaw == nil {
		return nil, errors.From(errors.New("signer role not found"), logan.F{
			"role_id": signer.RoleID,
		})
	}

	role := resources.NewSignerRole(*roleRaw)

	rules, err := h.SignerRuleQ.FilterByRole(roleRaw.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load signer role rules for role")
	}

	for _, rawRule := range rules {
		rule := resources.NewSignerRule(rawRule)
		role.Relationships.Rules.Data = append(role.Relationships.Rules.Data, rule.Key)
		if request.ShouldInclude(requests.IncludeTypeSignerRoleRules) {
			includes.Add(&rule)
		}
	}

	includes.Add(&role)

	return role.AsRelation(), nil
}
