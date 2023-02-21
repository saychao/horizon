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

// GetPublicKeyEntry - processes request to get the public key with the list of all accounts that have signer with provided public key
func GetPublicKeyEntry(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getPublicKeyEntryHandler{
		SignersQ: core2.NewSignerQ(coreRepo),
		Log:      ctx.Log(r),
	}

	request, err := requests.NewGetPublicKeyEntry(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, &request.ID) {
		return
	}

	result, err := handler.GetPublicKeyEntry(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get public key entry", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if result == nil {
		ape.RenderErr(w, problems.NotFound())
	}

	ape.Render(w, result)
}

type getPublicKeyEntryHandler struct {
	SignersQ core2.SignerQ
	Log      *logan.Entry
}

// GetPublicKeyEntry returns the public key entry with the list of related account keys
func (h *getPublicKeyEntryHandler) GetPublicKeyEntry(request *requests.GetPublicKeyEntry) (*regources.PublicKeyEntryResponse, error) {
	publicKeyEntry := resources.NewPublicKeyEntry(request.ID)

	coreSigners, err := h.SignersQ.FilterByPublicKey(request.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get signers list")
	}
	if len(coreSigners) == 0 {
		return nil, nil
	}

	accounts := make([]regources.Key, 0, len(coreSigners))
	for _, coreSigner := range coreSigners {
		accounts = append(accounts, resources.NewAccountKey(coreSigner.AccountID))
	}

	publicKeyEntry.Relationships.Accounts = &regources.RelationCollection{
		Data: accounts,
	}

	response := &regources.PublicKeyEntryResponse{
		Data: publicKeyEntry,
	}

	return response, nil
}
