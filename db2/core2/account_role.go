package core2

import regources "gitlab.com/tokend/regources/generated"

//AccountRole - represents role applicable for account
type AccountRole struct {
	ID      uint64            `db:"id"`
	Details regources.Details `db:"details"`
}
