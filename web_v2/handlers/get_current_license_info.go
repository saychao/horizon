package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetCurrentLicenseInfo(w http.ResponseWriter, r *http.Request) {
	handler := getLicenseHandler{
		LicensesQ: core2.NewLicenseQ(ctx.CoreRepo(r)),
		SignersQ:  core2.NewSignerQ(ctx.CoreRepo(r)),
		AdminID:   ctx.CoreInfo(r).AdminAccountID,
		Log:       ctx.Log(r),
	}

	result, err := handler.getLicense()
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get license")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !isAllowed(r, w) {
		return
	}

	ape.Render(w, result)
}

type getLicenseHandler struct {
	LicensesQ core2.LicenseQ
	SignersQ  core2.SignerQ
	AdminID   string
	Log       *logan.Entry
}

// GetSale returns sale with related resources
func (h *getLicenseHandler) getLicense() (regources.LicenseInfoResponse, error) {

	record, err := h.LicensesQ.GetLatest()
	if err != nil {
		return regources.LicenseInfoResponse{}, errors.Wrap(err, "failed to get license")
	}

	if record == nil {
		return regources.LicenseInfoResponse{}, nil
	}

	admins, err := h.SignersQ.Count(h.AdminID)
	if err != nil {
		return regources.LicenseInfoResponse{}, errors.Wrap(err, "failed to get number of admins")
	}

	resp := regources.LicenseInfoResponse{}

	resp.Data = resources.NewLicenseInfo(*record, admins)

	return resp, nil
}
