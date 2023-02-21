package resources

import (
	regources "gitlab.com/tokend/regources/generated"
)

// NewPublicKeyEntry creates new instance of PublicKeyEntry resource from provided publicKey
func NewPublicKeyEntry(publicKey string) regources.PublicKeyEntry {
	return regources.PublicKeyEntry{
		Key: regources.Key{
			ID:   publicKey,
			Type: regources.PUBLIC_KEY_ENTRIES,
		},
	}
}
