package core2

import regources "gitlab.com/tokend/regources/generated"

type AccountKYC struct {
	AccountID string            `db:"accountid"`
	KYCData   regources.Details `db:"kyc_data"`
}
