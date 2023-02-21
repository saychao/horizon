package handlers

import (
	"net/http"

	"github.com/saychao/horizon/db2/history2"

	"github.com/saychao/horizon/db2/core2"
	fees "github.com/saychao/horizon/helper/fees2"
	"github.com/saychao/horizon/helper/limits"
	"github.com/saychao/horizon/helper/statslimits"
	"github.com/saychao/horizon/web_v2/ctx"
	"github.com/saychao/horizon/web_v2/requests"
	"github.com/saychao/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// GetAccount - processes request to get account and it's details by address
func GetAccount(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	handler := getAccountHandler{
		AccountsQ:          core2.NewAccountsQ(coreRepo),
		HistoryAccountsQ:   history2.NewAccountsQ(ctx.HistoryRepo(r)),
		BalancesQ:          core2.NewBalancesQ(coreRepo),
		AccountRoleQ:       core2.NewAccountRoleQ(coreRepo),
		AccountRuleQ:       core2.NewAccountRuleQ(coreRepo),
		FeesQ:              core2.NewFeesQ(coreRepo),
		LimitsV2Q:          core2.NewLimitsQ(coreRepo),
		ExternalSystemIDsQ: core2.NewExternalSystemIDsQ(coreRepo),
		StatsQ:             core2.NewStatsQ(coreRepo),
		KycQ:               core2.NewAccountsKycQ(coreRepo),
		Log:                ctx.Log(r),
	}

	request, err := requests.NewGetAccount(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO: MUST be changes to role based access
	if request.ShouldIncludeAny(
		requests.IncludeTypeAccountBalancesState,
		requests.IncludeTypeAccountAccountReferrer,
		requests.IncludeTypeAccountFees,
		requests.IncludeTypeAccountLimits,
		requests.IncludeTypeAccountExternalSystemIDs,
		requests.IncludeTypeAccountLimitsWithStats,
		requests.IncludeTypeAccountKycData,
	) {
		if !isAllowed(r, w, &request.Address) {
			return
		}
	}

	result, err := handler.GetAccount(request)
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

type getAccountHandler struct {
	AccountsQ          core2.AccountsQ
	BalancesQ          core2.BalancesQ
	AccountRoleQ       core2.AccountRoleQ
	AccountRuleQ       core2.AccountRuleQ
	LimitsV2Q          core2.LimitsQ
	FeesQ              core2.FeesQ
	ExternalSystemIDsQ core2.ExternalSystemIDsQ
	StatsQ             core2.StatsQ
	KycQ               core2.AccountsKycQ
	HistoryAccountsQ   history2.AccountsQ
	Log                *logan.Entry
}

// GetAccount - returns Account resources
func (h *getAccountHandler) GetAccount(request *requests.GetAccount) (*regources.AccountResponse, error) {
	account, err := h.AccountsQ.GetByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by address")
	}

	if account == nil {
		return nil, nil
	}

	response := regources.AccountResponse{
		Data: resources.NewAccount(*account, nil),
	}

	accountStatus, err := h.HistoryAccountsQ.ByAddress(request.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account status")
	}

	if accountStatus != nil {
		recoveryStatus := regources.KYCRecoveryStatus(accountStatus.KycRecoveryStatus)
		response.Data.Attributes.KycRecoveryStatus = &recoveryStatus
	}

	response.Data.Relationships.Role, err = h.getRole(request, &response.Included, *account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get role")
	}

	response.Data.Relationships.Balances, err = h.getBalances(request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get balances")
	}

	response.Data.Relationships.Referrer, err = h.getReferrer(account, request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get referrer")
	}

	response.Data.Relationships.Fees, err = h.getFees(request, &response.Included, account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees for account")
	}

	response.Data.Relationships.Limits, err = h.getLimits(request, &response.Included, account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get limits")
	}

	response.Data.Relationships.ExternalSystemIds, err = h.getExternalSystemIDs(request, &response.Included)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get external system IDs")
	}

	response.Data.Relationships.LimitsWithStats, err = h.getLimitsWithStats(request, &response.Included, account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get limits and stats for account")
	}

	response.Data.Relationships.KycData, err = h.getKycData(request, &response.Included, account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kyc data for account")
	}

	return &response, nil
}

func (h *getAccountHandler) getLimits(request *requests.GetAccount,
	includes *regources.Included,
	account *core2.Account) (*regources.RelationCollection, error) {

	generalLimits, err := h.LimitsV2Q.General().Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select general limits")
	}
	roleLimits, err := h.LimitsV2Q.FilterByAccountRole(account.RoleID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select role limits")
	}
	accountLimits, err := h.LimitsV2Q.FilterByAccount(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select limits for account")
	}

	lt := limits.NewTable(generalLimits)
	lt.Update(roleLimits)
	lt.Update(accountLimits)

	result := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(lt)),
	}

	for _, coreLimitsUnit := range lt {
		limitsUnit := resources.NewLimits(coreLimitsUnit)
		result.Data = append(result.Data, limitsUnit.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountLimits) {
			includes.Add(&limitsUnit)
		}
	}

	return result, nil
}

func (h *getAccountHandler) getLimitsWithStats(request *requests.GetAccount,
	included *regources.Included,
	account *core2.Account) (*regources.RelationCollection, error) {

	generalLimits, err := h.LimitsV2Q.General().Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select general limits")
	}
	roleLimits, err := h.LimitsV2Q.FilterByAccountRole(account.RoleID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select role limits")
	}
	accountLimits, err := h.LimitsV2Q.FilterByAccount(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select limits for account")
	}
	accountStats, err := h.StatsQ.FilterByAccount(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select statistics for account")
	}

	statsLimitsTable := statslimits.NewTable(generalLimits, accountStats)
	statsLimitsTable.Update(roleLimits)
	statsLimitsTable.Update(accountLimits)
	statsLimitsTable.FulfillEmptyLimits()

	coreLimitsWithStatsList := statsLimitsTable.CoreUnitsList()

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreLimitsWithStatsList)),
	}

	for _, coreLimitsWithStatsUnit := range coreLimitsWithStatsList {
		limitsWithStatsUnit := resources.NewLimitsWithStats(&coreLimitsWithStatsUnit)

		result.Data = append(result.Data, limitsWithStatsUnit.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountLimitsWithStats) {
			included.Add(&limitsWithStatsUnit)
		}
	}
	return &result, nil
}

func (h *getAccountHandler) getExternalSystemIDs(request *requests.GetAccount, includes *regources.Included) (*regources.RelationCollection, error) {
	externalSystemIDsQ := h.ExternalSystemIDsQ.FilterByAccount(request.Address)

	coreExternalSystemIDs, err := externalSystemIDsQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select external system IDs for account")
	}

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreExternalSystemIDs)),
	}

	for _, coreExtSysIDUnit := range coreExternalSystemIDs {
		externalSystemID := resources.NewExternalSystemID(coreExtSysIDUnit)

		result.Data = append(result.Data, externalSystemID.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountExternalSystemIDs) {
			includes.Add(externalSystemID)
		}
	}
	return &result, nil
}

func (h *getAccountHandler) getReferrer(account *core2.Account, request *requests.GetAccount, includes *regources.Included) (*regources.Relation, error) {
	if account.Referrer == nil {
		return nil, nil
	}

	if !request.ShouldInclude(requests.IncludeTypeAccountAccountReferrer) {
		result := resources.NewAccountKey(*account.Referrer)

		return result.AsRelation(), nil
	}

	referrer, err := h.AccountsQ.GetByAddress(*account.Referrer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load referrer", logan.F{
			"referrer_address": account.Referrer,
		})
	}

	if referrer == nil {
		return nil, errors.Wrap(err, "expected referrer to exist", logan.F{
			"account_address":  account.Address,
			"referrer_address": account.Referrer,
		})
	}

	result := resources.NewAccount(*referrer, nil)
	includes.Add(&result)

	return result.AsRelation(), nil
}

func (h *getAccountHandler) getBalances(request *requests.GetAccount, includes *regources.Included) (*regources.RelationCollection, error) {
	balancesQ := h.BalancesQ.FilterByAccount(request.Address)

	if request.ShouldInclude(requests.IncludeTypeAccountBalancesAsset) {
		balancesQ = balancesQ.WithAsset()
	}

	coreBalances, err := balancesQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select balances for account")
	}

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0, len(coreBalances)),
	}

	for i, coreBalance := range coreBalances {
		balance := resources.NewBalance(&coreBalances[i])
		result.Data = append(result.Data, balance.Key)

		if request.ShouldInclude(requests.IncludeTypeAccountBalances) {
			balance.Relationships = &regources.BalanceRelationships{}
			balance.Relationships.State = resources.NewBalanceStateKey(coreBalance.BalanceAddress).AsRelation()
			balance.Relationships.Asset = resources.NewAssetKey(coreBalance.AssetCode).AsRelation()

			includes.Add(balance)
		}

		if request.ShouldInclude(requests.IncludeTypeAccountBalancesState) {
			state := resources.NewBalanceState(&coreBalances[i])
			includes.Add(state)
		}

		if request.ShouldInclude(requests.IncludeTypeAccountBalancesAsset) {
			asset := resources.NewAsset(*coreBalances[i].Asset)
			includes.Add(&asset)
		}
	}

	return &result, nil
}

func (h *getAccountHandler) getRole(request *requests.GetAccount,
	includes *regources.Included, account core2.Account,
) (*regources.Relation, error) {
	if !request.ShouldInclude(requests.IncludeTypeAccountRole) {
		role := resources.NewAccountRoleKey(account.RoleID)
		return role.AsRelation(), nil
	}

	roleRaw, err := h.AccountRoleQ.FilterByID(account.RoleID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account role by id")
	}

	if roleRaw == nil {
		return nil, errors.From(errors.New("role not found"), logan.F{
			"id": account.RoleID,
		})
	}

	role := resources.NewAccountRole(*roleRaw)
	rules, err := h.AccountRuleQ.FilterByRole(roleRaw.ID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select account rules for role")
	}

	for _, ruleRaw := range rules {
		rule := resources.NewAccountRule(ruleRaw)
		role.Relationships.Rules.Data = append(role.Relationships.Rules.Data, rule.GetKey())
		if request.ShouldInclude(requests.IncludeTypeAccountRoleRules) {
			includes.Add(&rule)
		}
	}

	return role.AsRelation(), nil
}

func (h *getAccountHandler) getFees(request *requests.GetAccount, includes *regources.Included, account *core2.Account) (*regources.RelationCollection, error) {

	coreBalances, err := h.BalancesQ.FilterByAccount(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select balances for account")
	}

	assets := make([]string, 0, len(coreBalances))
	for _, coreBalance := range coreBalances {
		assets = append(assets, coreBalance.AssetCode)
	}

	forAccount, err := h.FeesQ.FilterByAddress(request.Address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees for account")
	}

	forAccountRole, err := h.FeesQ.FilterByAccountRole(account.RoleID).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fees for account role")
	}

	generalFees, err := h.FeesQ.FilterGlobal().Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get global fees")
	}

	sft := fees.NewSmartFeeTable(forAccount)
	sft.Update(forAccountRole)
	sft.Update(generalFees)
	sft.AddZeroFees(assets)

	result := regources.RelationCollection{
		Data: make([]regources.Key, 0),
	}
	for _, v := range sft {
		for _, feeRecord := range v {
			hash := resources.CalculateFeeHash(feeRecord.Fee)
			feeKey := resources.NewFeeKey(hash)
			result.Data = append(result.Data, feeKey)

			if request.ShouldInclude(requests.IncludeTypeAccountFees) {
				fee := resources.NewFee(feeRecord.Fee)
				includes.Add(&fee)
			}
		}
	}

	return &result, nil
}

func (h *getAccountHandler) getKycData(request *requests.GetAccount, includes *regources.Included, account *core2.Account) (*regources.Relation, error) {
	kycDataRaw, err := h.KycQ.GetByAddress(account.Address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kyc data for account", logan.F{
			"address": account.Address,
		})
	}
	if kycDataRaw == nil {
		return nil, nil
	}
	kycData := resources.NewAccountKYC(*kycDataRaw)

	if request.ShouldInclude(requests.IncludeTypeAccountKycData) {
		includes.Add(&kycData)
	}

	return kycData.AsRelation(), nil
}
