package resources

import (
	"strconv"

	"github.com/saychao/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewSigner - creates new instance of signer
func NewSigner(signer core2.Signer) regources.Signer {
	return regources.Signer{
		Key: NewSignerKey(signer.PublicKey),
		Attributes: regources.SignerAttributes{
			Weight:   signer.Weight,
			Identity: signer.Identity,
			Details:  signer.Details,
		},
		Relationships: regources.SignerRelationships{
			Role: regources.Key{
				ID:   strconv.FormatUint(signer.RoleID, 10),
				Type: regources.SIGNER_ROLES,
			}.AsRelation(),
			Account: regources.Key{
				ID:   signer.AccountID,
				Type: regources.ACCOUNTS,
			}.AsRelation(),
		},
	}
}

// NewSignerKey - creates new key for signer
func NewSignerKey(publicKey string) regources.Key {
	return regources.Key{
		ID:   publicKey,
		Type: regources.SIGNERS,
	}
}
