package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	regources "gitlab.com/tokend/regources/generated"
)

const maximumTrailingDigits uint32 = 6

// GetCalculatedFees - processes request to get the list of fees
func GetCalculatedFees(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getCalculatedFeesHandler{
		FeesQ:     core2.NewFeesQ(coreRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		AssetsQ:   core2.NewAssetsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetCalculatedFees(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetCalculatedFees(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to calculate fee", logan.F{
			"address":  request.Address,
			"amount":   request.Amount,
			"fee_type": request.FeeType,
			"subtype":  request.Subtype,
			"asset":    request.Asset,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getCalculatedFeesHandler struct {
	FeesQ     core2.FeesQ
	AccountsQ core2.AccountsQ
	AssetsQ   core2.AssetsQ
	Log       *logan.Entry
}

// GetCalculatedFees returns calculated fee for given given parameters
func (h *getCalculatedFeesHandler) GetCalculatedFees(request *requests.GetCalculatedFees) (*regources.CalculatedFeeResponse, error) {
	fee, err := h.getFeeForAccount(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load fee for account")
	}

	minimalAmount, err := h.getMinimalAmount(fee.Asset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate minimal amount for asset")
	}

	calculatedPercent, overflow := amount.BigDivide(int64(request.Amount), fee.Percent, 100*amount.One, amount.ROUND_UP, minimalAmount)
	if overflow {
		return nil, errors.New("failed to calculate fee")
	}

	attributes := regources.Fee{
		CalculatedPercent: regources.Amount(calculatedPercent),
		Fixed:             regources.Amount(fee.Fixed),
	}

	hash := h.getHash(attributes)
	response := &regources.CalculatedFeeResponse{
		Data: regources.CalculatedFee{
			Key:        resources.NewCalculatedFeeKey(hash),
			Attributes: attributes,
		},
	}

	return response, nil
}

func (h *getCalculatedFeesHandler) getMinimalAmount(assetCode string) (int64, error) {
	asset, err := h.AssetsQ.GetByCode(assetCode)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get asset")
	}
	if asset == nil {
		return 0, errors.New("asset not found")
	}

	exp := maximumTrailingDigits - asset.TrailingDigits
	if exp < 0 {
		return 0, errors.New("Incorrect asset precision")
	}
	minimalAmount := int64(1)
	for exp > 0 {
		minimalAmount *= 10
		exp -= 1
	}

	return minimalAmount, nil
}

func (h *getCalculatedFeesHandler) getHash(attrs regources.Fee) string {
	data := fmt.Sprintf("fixed:%s:calculated_percent:%s", attrs.Fixed.String(), attrs.CalculatedPercent.String())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (h *getCalculatedFeesHandler) getFeeForAccount(request *requests.GetCalculatedFees) (*core2.Fee, error) {
	q := h.FeesQ.
		FilterByAsset(request.Asset).
		FilterBySubtype(request.Subtype).
		FilterByType(request.FeeType).
		FilterByAmount(int64(request.Amount))

	//try get fee for account
	fee, err := q.FilterByAddress(request.Address).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fee for address")
	}

	if fee != nil {
		return fee, nil
	}

	targetAccount, err := h.AccountsQ.FilterByAddress(request.Address).Get()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get account")
	}
	if targetAccount == nil {
		return nil, errors.New("Account not found")
	}
	//try to get fee for account type
	fee, err = q.FilterByAccountRole(targetAccount.RoleID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fee for account type")
	}

	if fee != nil {
		return fee, nil
	}

	fee, err = q.FilterGlobal().Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load general fees")
	}
	if fee != nil {
		return fee, nil
	}

	return &core2.Fee{
		AccountRole: core2.FeesEmptyRole,
		FeeType:     request.FeeType,
		Subtype:     request.Subtype,
		Asset:       request.Asset,
		Fixed:       0,
		Percent:     0,
	}, nil

}
