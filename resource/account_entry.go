package resource

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type AccountEntry struct {
	AccountID     string            `json:"account_id"`
	AccountTypeI  int32             `json:"account_type_i"`
	AccountType   string            `json:"account_type"`
	BlockReasonsI uint32            `json:"block_reasons_i"`
	BlockReasons  []regources.Flag  `json:"block_reasons"`
	LimitsV2      []LimitsV2        `json:"limits"`
	Policies      AccountPolicies   `json:"policies"`
	Thresholds    AccountThresholds `json:"thresholds"`
}

type LedgerKeyAccount struct {
	AccountID string `json:"account_id"`
}

func (r *AccountEntry) Populate(entry xdr.AccountEntry) {
	r.AccountID = entry.AccountId.Address()
}

func (r *LedgerKeyAccount) Populate(xdrAcc xdr.LedgerKeyAccount) {
	r.AccountID = xdrAcc.AccountId.Address()
}
