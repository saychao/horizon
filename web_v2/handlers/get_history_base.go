package handlers

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/saychao/horizon/db2/core2"

	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

type getHistory struct {
	// cached filter entries
	Account *history2.Account
	Balance *history2.Balance
	Asset   *core2.Asset

	BalancesIDs []uint64

	AssetsQ   core2.AssetsQ
	EffectsQ  history2.ParticipantEffectsQ
	AccountsQ history2.AccountsQ
	BalanceQ  history2.BalancesQ
	Log       *logan.Entry
}

func newHistoryHandler(r *http.Request) getHistory {
	historyRepo := ctx.HistoryRepo(r)
	handler := getHistory{
		AssetsQ:   core2.NewAssetsQ(ctx.CoreRepo(r)),
		EffectsQ:  history2.NewParticipantEffectsQ(historyRepo),
		AccountsQ: history2.NewAccountsQ(historyRepo),
		BalanceQ:  history2.NewBalancesQ(historyRepo),
		Log:       ctx.Log(r),
	}

	return handler
}

func (h *getHistory) prepare(w http.ResponseWriter, r *http.Request) (*requests.GetHistory, bool) {
	request, err := requests.NewGetHistory(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return nil, false
	}

	// TODO: need to refactor
	if request.Filters.Account != nil {
		h.Account, err = h.AccountsQ.ByAddress(*request.Filters.Account)
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to get account", logan.F{
				"account_address": request.Filters.Account,
			})
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}

		if h.Account == nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"filter[account]": errors.New("not found"),
			})...)
			return nil, false
		}
	}

	if len(request.Filters.Balance) > 1 {
		balances, err := h.BalanceQ.SelectByAddress(request.Filters.Balance...)
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to get balances", logan.F{
				"balances": request.Filters.Balance,
			})
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}

		loadedBalances := make(map[string]bool)
		differentBalances := make([]string, 0, len(balances))
		for idx, balance := range balances {
			loadedBalances[balance.Address] = true
			h.BalancesIDs = append(h.BalancesIDs, balances[idx].ID)
			if balance.AccountID != balances[0].AccountID {
				differentBalances = append(differentBalances, balance.Address)
			}
		}
		if len(differentBalances) != 0 {
			ctx.Log(r).Error("balances must belong to the same account", logan.F{
				"balance": differentBalances,
			})
			ape.RenderErr(w, problems.Conflict())
			return nil, false
		}

		if len(balances) != len(request.Filters.Balance) {
			notFoundBalances := make([]string, 0, len(balances))
			for _, balance := range request.Filters.Balance {
				if _, loaded := loadedBalances[balance]; !loaded {
					notFoundBalances = append(notFoundBalances, balance)
				}
			}

			if len(notFoundBalances) > 0 {
				ctx.Log(r).Error("balance not found", logan.F{
					"balance": notFoundBalances,
				})
				ape.RenderErr(w, problems.NotFound())
			} else {
				ctx.Log(r).Error("balance duplicate")
				ape.RenderErr(w, problems.Conflict())
			}

			return nil, false
		}

		if !h.ensureAllowed(w, r) {
			return nil, false
		}

		return request, true
	}

	if len(request.Filters.Balance) == 1 {
		h.Balance, err = h.BalanceQ.GetByAddress(request.Filters.Balance[0])
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to get balance", logan.F{
				"balance_address": request.Filters.Balance,
			})
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}

		if h.Balance == nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"filter[balance]": errors.New("not found"),
			})...)
			return nil, false
		}

		h.BalancesIDs = append(h.BalancesIDs, h.Balance.ID)
	}

	if request.Filters.Asset != nil {
		h.Asset, err = h.AssetsQ.GetByCode(*request.Filters.Asset)
		if err != nil {
			ctx.Log(r).WithError(err).Error("failed to get asset", logan.F{
				"asset_code": request.Filters.Asset,
			})
			ape.RenderErr(w, problems.InternalError())
			return nil, false
		}

		if h.Asset == nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"filter[asset]": errors.New("not found"),
			})...)
			return nil, false
		}
	}

	if !h.ensureAllowed(w, r) {
		return nil, false
	}

	return request, true
}

func (h *getHistory) ApplyFilters(request *requests.GetHistory,
	q history2.ParticipantEffectsQ) history2.ParticipantEffectsQ {
	q = q.WithAccount().WithBalance().Page(request.PageParams)
	if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
		q = q.WithOperation()
	}

	if h.Account != nil {
		q = q.ForAccount(h.Account.ID)
	}

	if len(h.BalancesIDs) > 0 {
		q = q.ForBalance(h.BalancesIDs...)
	}

	if h.Asset != nil {
		q = q.ForAsset(h.Asset.Code)
	}

	if request.Filters.OpsCloseTimeFrom != nil {
		q = q.LedgerCloseTimeFrom(*request.Filters.OpsCloseTimeFrom)
	}

	if request.Filters.OpsCloseTimeTo != nil {
		q = q.LedgerCloseTimeTo(*request.Filters.OpsCloseTimeTo)
	}

	if len(request.Filters.EffectType) > 0 {
		effects := make([]history2.EffectType, len(request.Filters.EffectType))
		for i, effect := range request.Filters.EffectType {
			effects[i] = resources.EffectTypeFromString(regources.ResourceType(effect))
		}
		q = q.ForEffect(effects...)
	}

	return q
}

func (h *getHistory) SelectAndPopulate(
	request *requests.GetHistory,
	effectsQ history2.ParticipantEffectsQ,
	assetQ core2.AssetsQ,
) (regources.ParticipantsEffectListResponse, error) {

	result := regources.ParticipantsEffectListResponse{
		Data: []regources.ParticipantsEffect{},
	}

	effects, err := effectsQ.Select()
	if err != nil {
		return result, errors.Wrap(err, "failed to load participant effects")
	}

	if len(effects) == 0 {
		result.Links = request.GetCursorLinks(request.PageParams, "")
		return result, nil
	}

	result.Data = make([]regources.ParticipantsEffect, 0, len(effects))
	for i := range effects {
		effect := getEffect(effects[i])
		if request.ShouldInclude(requests.IncludeTypeHistoryOperation) {
			op := resources.NewOperation(*effects[i].Operation)

			opDetails := resources.NewOperationDetails(*effects[i].Operation)
			op.Relationships.Details = opDetails.GetKey().AsRelation()

			if request.ShouldInclude(requests.IncludeTypeHistoryOperationDetails) {
				result.Included.Add(opDetails)
			}

			result.Included.Add(&op)
		}

		if effects[i].Effect != nil {
			change := resources.NewEffect(effects[i].ID, *effects[i].Effect)
			effect.Relationships.Effect = change.GetKey().AsRelation()
			if request.ShouldInclude(requests.IncludeTypeHistoryEffect) {
				result.Included.Add(change)
			}
		}

		if effects[i].AssetCode != nil {
			if request.ShouldInclude(requests.IncludeTypeHistoryAsset) {
				rawAsset, err := assetQ.GetByCode(*effects[i].AssetCode)
				if err != nil {
					return result, errors.Wrap(err, "failed to load asset")
				}
				asset := resources.NewAsset(*rawAsset)
				result.Included.Add(&asset)
			}
		}

		result.Data = append(result.Data, effect)
	}

	result.Links = request.GetCursorLinks(request.PageParams, result.Data[len(result.Data)-1].ID)

	return result, nil
}

func getEffect(effect history2.ParticipantEffect) regources.ParticipantsEffect {
	var balance *regources.Relation
	if effect.BalanceAddress != nil {
		balance = resources.NewBalanceKey(*effect.BalanceAddress).AsRelation()
	}

	var asset *regources.Relation
	if effect.AssetCode != nil {
		asset = resources.NewAssetKey(*effect.AssetCode).AsRelation()
	}

	return regources.ParticipantsEffect{
		Key: resources.NewParticipantEffectKey(effect.ID),
		Relationships: regources.ParticipantsEffectRelationships{
			Account:   resources.NewAccountKey(effect.AccountAddress).AsRelation(),
			Balance:   balance,
			Asset:     asset,
			Operation: resources.NewOperationKey(effect.OperationID).AsRelation(),
		},
	}

}

// ensure allowed - checks it requester is allowed to access the data. If not it renders error and returns false.
// The logic behind this is that if multiple filters provided all resource owners have access to data, as we
// returning smaller subset of effects/movements
func (h *getHistory) ensureAllowed(w http.ResponseWriter, httpRequest *http.Request) bool {
	constraints := make([]*string, 0)
	if h.Account != nil {
		constraints = append(constraints, &h.Account.Address)
	}

	if h.Balance != nil {
		account, err := h.tryGetAccountForBalance(h.Balance)
		if err != nil {
			ctx.Log(httpRequest).WithError(err).Error("failed to load account for balance")
			ape.RenderErr(w, problems.InternalError())
			return false
		}

		constraints = append(constraints, &account.Address)
	}

	if h.Asset != nil {
		constraints = append(constraints, &h.Asset.Owner)
	}
	// Admin is added implicitly to constraints in `isAllowed`, so no need to add it explicitly
	return isAllowed(httpRequest, w, constraints...)
}

func (h *getHistory) tryGetAccountForBalance(balance *history2.Balance) (history2.Account, error) {
	// TODO: fuck normalization
	account, err := h.AccountsQ.ByID(balance.AccountID)
	if err != nil {
		return history2.Account{}, errors.Wrap(err, "failed to load account by ID")
	}

	if account == nil {
		return history2.Account{}, errors.From(errors.New("found balance, but failed to find account for it"), logan.F{
			"balance": balance.ID,
		})
	}

	return *account, nil
}
