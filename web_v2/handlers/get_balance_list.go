package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
)

// GetBalanceList - processes request to get the list of balances
func GetBalanceList(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getBalanceListHandler{
		AssetsQ:          core2.NewAssetsQ(coreRepo),
		BalancesQ:        core2.NewBalancesQ(coreRepo),
		AccountsQ:        core2.NewAccountsQ(coreRepo),
		HistoryAccountsQ: history2.NewAccountsQ(ctx.HistoryRepo(r)),
		Log:              ctx.Log(r),
	}

	request, err := requests.NewGetBalanceList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var assetOwner string
	if request.Filters.Asset != nil {
		assetOwner, err = handler.getAssetOwner(*request.Filters.Asset)

		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to get asset owner", logan.F{
				"request": request,
			})
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}
	if !isAllowed(r, w, &assetOwner, request.Filters.Owner) {
		return
	}

	result, err := handler.GetBalanceList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get asset list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getBalanceListHandler struct {
	BalancesQ        core2.BalancesQ
	AssetsQ          core2.AssetsQ
	AccountsQ        core2.AccountsQ
	HistoryAccountsQ history2.AccountsQ
	Log              *logan.Entry
}

func (h *getBalanceListHandler) getAssetOwner(assetCode string) (string, error) {
	if assetCode == "" {
		return "", nil
	}

	coreAsset, err := h.AssetsQ.GetByCode(assetCode)
	if err != nil {
		return "", errors.Wrap(err, "failed to get asset")
	}

	if coreAsset == nil {
		return "", nil
	}

	return coreAsset.Owner, nil
}

// GetBalanceList returns list of balances with related resources
func (h *getBalanceListHandler) GetBalanceList(request *requests.GetBalanceList) (*regources.BalanceListResponse, error) {
	q := h.BalancesQ.Page(request.PageParams)
	if request.Filters.Asset != nil {
		q = q.FilterByAsset(*request.Filters.Asset)
	}

	if request.Filters.AssetOwner != nil {
		q = q.FilterByAssetOwner(*request.Filters.AssetOwner)
	}

	if request.Filters.Owner != nil {
		q = q.FilterByAccount(*request.Filters.Owner)
	}

	if request.Filters.AmountGt != nil {
		q = q.FilterByAmountGt(*request.Filters.AmountGt)
	}

	if request.Filters.AmountLwOrEq != nil {
		q = q.FilterByAmountLwOrEq(*request.Filters.AmountLwOrEq)
	}

	coreBalances, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get balance list")
	}

	response := &regources.BalanceListResponse{
		Data:  make([]regources.Balance, 0, len(coreBalances)),
		Links: request.GetOffsetLinks(request.PageParams),
	}

	for _, coreBalance := range coreBalances {
		balance := resources.NewBalance(&coreBalance)
		balance.Relationships = &regources.BalanceRelationships{
			Asset: resources.NewAssetKey(coreBalance.AssetCode).AsRelation(),
			State: resources.NewBalanceStateKey(coreBalance.BalanceAddress).AsRelation(),
			Owner: resources.NewAccountKey(coreBalance.AccountAddress).AsRelation(),
		}

		if request.ShouldInclude(requests.IncludeTypeBalanceListState) {
			response.Included.Add(resources.NewBalanceState(&coreBalance))
		}

		if request.ShouldInclude(requests.IncludeTypeBalanceListOwner) {
			account, err := h.AccountsQ.GetByAddress(coreBalance.AccountAddress)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get account by address")
			}

			if account == nil {
				return nil, errors.New("owner not found")
			}

			accountStatus, err := h.HistoryAccountsQ.ByAddress(coreBalance.AccountAddress)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get account status")
			}
			if accountStatus == nil {
				return nil, errors.New("account not found in history")
			}
			recoveryStatus := regources.KYCRecoveryStatus(accountStatus.KycRecoveryStatus)
			owner := resources.NewAccount(*account, &recoveryStatus)
			response.Included.Add(&owner)
		}

		response.Data = append(response.Data, *balance)
	}

	return response, nil
}
