package resources

import (
	"github.com/SafeRE-IT/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"
)

// NewAccountKYC - creates new account KYC resource
func NewAccountKYC(kyc core2.AccountKYC) regources.AccountKyc {
	return regources.AccountKyc{
		Key: NewAccountKYCKey(kyc.AccountID),
		Attributes: regources.AccountKycAttributes{
			KycData: kyc.KYCData,
		},
	}
}

//NewAccountRoleKey - returns new instance of key for account role
func NewAccountKYCKey(address string) regources.Key {
	return regources.Key{
		ID:   address,
		Type: regources.ACCOUNT_KYC,
	}
}
