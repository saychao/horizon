package resource

import (
	"github.com/saychao/horizon/db2/core"
)

type Referral struct {
	AccountID string `json:"account_id"`
	RoleID    uint64 `json:"role_id"`
}

// Populate fills out the resource's fields
func (r *Referral) Populate(ca core.Account) {
	r.AccountID = ca.AccountID
	r.RoleID = ca.RoleID
}
