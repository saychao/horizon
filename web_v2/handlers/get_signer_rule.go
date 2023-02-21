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

// GetSignerRule - processes request to get signerRule and it's details by signerRule code
func GetSignerRule(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getSignerRuleHandler{
		SignerRulesQ: core2.NewSignerRuleQ(coreRepo),
		Log:          ctx.Log(r),
	}

	request, err := requests.NewGetSignerRule(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSignerRule(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get signer rule", logan.F{
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

type getSignerRuleHandler struct {
	SignerRulesQ core2.SignerRuleQ
	Log          *logan.Entry
}

// GetSignerRule returns signerRule with related resources
func (h *getSignerRuleHandler) GetSignerRule(request *requests.GetSignerRule) (*regources.SignerRuleResponse, error) {
	signerRule, err := h.SignerRulesQ.FilterByID(request.ID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get signer rule by id")
	}
	if signerRule == nil {
		return nil, nil
	}

	signerRuleResponse := resources.NewSignerRule(*signerRule)
	return &regources.SignerRuleResponse{
		Data: signerRuleResponse,
	}, nil
}
