package resource

type AccountPolicyType struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

type AccountPolicies struct {
	AccountPoliciesTypeI int32               `json:"account_policies_type_i"`
	AccountPoliciesTypes []AccountPolicyType `json:"account_policies_types"`
}

func (ap *AccountPolicies) Populate(accountPoliciesType int32) {
	ap.AccountPoliciesTypeI = accountPoliciesType
}
