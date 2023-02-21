package core2

import (
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type SignerRule struct {
	ID        uint64                 `db:"id"`
	Resource  xdr.SignerRuleResource `db:"resource"`
	Action    int32                  `db:"action"`
	Forbids   bool                   `db:"forbids"`
	IsDefault bool                   `db:"is_default"`
	OwnerID   string                 `db:"owner_id"`
	Details   regources.Details      `db:"details"`
}
