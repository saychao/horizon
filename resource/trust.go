package resource

import (
	"github.com/SafeRE-IT/horizon/db2/core"
)

type Trust struct {
	AllowedAccount string `json:"allowed_account"`
	BalanceToUse   string `json:"balance_to_use"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (t *Trust) Populate(row core.Trust) {
	t.AllowedAccount = row.AllowedAccount
	t.BalanceToUse = row.BalanceToUse
}
