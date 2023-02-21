package requests

import (
	"net/http"
)

const (
	// IncludeTypeAccountBalances - defines if the account balances should be included in the response
	IncludeTypeAccountBalances = "balances"
	// IncludeTypeAccountBalancesAsset - defines if assets of the account balances should be included in the response
	IncludeTypeAccountBalancesAsset = "balances.asset"
	// IncludeTypeAccountBalancesState - defines if the account balances state should be included in the response
	IncludeTypeAccountBalancesState = "balances.state"
	// IncludeTypeAccountAccountReferrer - defines if the account referrer should be included in the response
	IncludeTypeAccountAccountReferrer = "referrer"
	// IncludeTypeAccountRole - defines if the account role should be included in the response
	IncludeTypeAccountRole = "role"
	// IncludeTypeAccountRoleRules - defines if rules of the account role should be included in the response
	IncludeTypeAccountRoleRules = "role.rules"
	//IncludeTypeAccountFees - defines if fees for the account should be included in the response
	IncludeTypeAccountFees = "fees"
	// IncludeTypeAccountLimits - defines if account limits should be included in the response
	IncludeTypeAccountLimits = "limits"
	// IncludeTypeAccountExternalSystemIDs - defines if account external system IDs should be included in the response
	IncludeTypeAccountExternalSystemIDs = "external_system_ids"
	// IncludeTypeAccountLimits - defines if account limits and statistics should be included in the response
	IncludeTypeAccountLimitsWithStats = "limits_with_stats"
	//IncludeTypeAccountKycData - defines if kyc data should be included in the response
	IncludeTypeAccountKycData = "kyc_data"
)

var includeTypeAccountAll = map[string]struct{}{
	IncludeTypeAccountBalances:          {},
	IncludeTypeAccountBalancesAsset:     {},
	IncludeTypeAccountBalancesState:     {},
	IncludeTypeAccountAccountReferrer:   {},
	IncludeTypeAccountRole:              {},
	IncludeTypeAccountRoleRules:         {},
	IncludeTypeAccountFees:              {},
	IncludeTypeAccountLimits:            {},
	IncludeTypeAccountExternalSystemIDs: {},
	IncludeTypeAccountLimitsWithStats:   {},
	IncludeTypeAccountKycData:           {},
}

//GetAccount - represents params to be specified by user for Get Account handler
type GetAccount struct {
	*base
	Address  string
	Includes struct {
		Balances          bool `include:"balances"`
		BalancesState     bool `include:"balances.state"`
		BalancesAsset     bool `include:"balances.asset"`
		Referrer          bool `include:"referrer"`
		Role              bool `include:"role"`
		RoleRules         bool `include:"role.rules"`
		Fees              bool `include:"fees"`
		Limits            bool `include:"limits"`
		ExternalSystemIds bool `include:"external_system_ids"`
		LimitsWithStats   bool `include:"limits_with_stats"`
		KycData           bool `include:"kyc_data"`
	}
}

//NewGetAccount - returns new instance of GetAccount request
func NewGetAccount(r *http.Request) (*GetAccount, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAccountAll,
	})

	if err != nil {
		return nil, err
	}
	address, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	return &GetAccount{
		base:    b,
		Address: address,
	}, nil

}
