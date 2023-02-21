package core2

import regources "gitlab.com/tokend/regources/generated"

//SignerRole - represents role of the signer
type SignerRole struct {
	ID      uint64            `db:"id"`
	RuleIDs []uint64          `db:"rule_ids"`
	OwnerID string            `db:"owner_id"`
	Details regources.Details `db:"details"`
}
