package resources

import (
	"strconv"

	"github.com/SafeRE-IT/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"
)

// NewAccountRole - creates role core account Role
func NewAccountRole(role core2.AccountRole) regources.AccountRole {
	return regources.AccountRole{
		Key: NewAccountRoleKey(role.ID),
		Attributes: regources.AccountRoleAttributes{
			Details: role.Details,
		},
		Relationships: regources.AccountRoleRelationships{
			Rules: regources.RelationCollection{},
		},
	}
}

// NewAccountRoleKey - returns new instance of key for account role
func NewAccountRoleKey(id uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: regources.ACCOUNT_ROLES,
	}
}

// NewSignerRole - maps signer role
func NewSignerRole(role core2.SignerRole) regources.SignerRole {
	return regources.SignerRole{
		Key: NewSignerRoleKey(role.ID),
		Attributes: regources.SignerRoleAttributes{
			Details: role.Details,
		},
		Relationships: regources.SignerRoleRelationships{
			Owner: NewAccountKey(role.OwnerID).AsRelation(),
			Rules: &regources.RelationCollection{},
		},
	}
}

// NewSignerRoleKey - creates new key for signer role
func NewSignerRoleKey(id uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: regources.SIGNER_ROLES,
	}
}
