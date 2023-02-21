package resource

import (
	"fmt"
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/httpx"
	"github.com/saychao/horizon/render/hal"
	"gitlab.com/tokend/regources"
	"golang.org/x/net/context"
)

// Account is the summary of an account
type Account struct {
	Links struct {
		Self         hal.Link `json:"self"`
		Transactions hal.Link `json:"transactions"`
		Operations   hal.Link `json:"operations"`
		Payments     hal.Link `json:"payments"`
	} `json:"_links"`

	HistoryAccount
	IsBlocked              bool              `json:"is_blocked"`
	BlockReasonsI          int32             `json:"block_reasons_i"`
	BlockReasons           []regources.Flag  `json:"block_reasons"`
	RoleID                 uint64            `json:"role_id"`
	AccountType            string            `json:"account_type"`
	Referrer               string            `json:"referrer"`
	Thresholds             AccountThresholds `json:"thresholds"`
	Balances               []Balance         `json:"balances"`
	Policies               AccountPolicies   `json:"policies"`
	AccountKYC             `json:"account_kyc"`
	ExternalSystemAccounts []regources.ExternalSystemAccountID `json:"external_system_accounts"`
	Referrals              []Referral                          `json:"referrals"`
}

// Populate fills out the resource's fields
func (a *Account) Populate(ctx context.Context, ca core.Account) {
	a.ID = ca.AccountID
	a.AccountID = ca.AccountID
	a.RoleID = ca.RoleID
	if ca.Referrer != nil {
		a.Referrer = *ca.Referrer
	}
	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	self := fmt.Sprintf("/accounts/%s", ca.AccountID)
	a.Links.Self = lb.Link(self)
	a.Links.Transactions = lb.PagedLink(self, "transactions")
	a.Links.Operations = lb.PagedLink(self, "operations")
	a.Links.Payments = lb.PagedLink(self, "payments")
	a.AccountKYC.Populate(*ca.AccountKYC)
}

func (a *Account) SetBalances(balances []core.Balance) {
	a.Balances = make([]Balance, len(balances))
	for i := range balances {
		a.Balances[i].Populate(balances[i])
	}
}

func (a Account) PagingToken() string {
	return a.ID
}
