package core

import "github.com/guregu/null"

// AccountKYC is a row of data from the `account_KYC` table
type AccountKYC struct {
	KYCData null.String `db:"account_kyc_data"`
}
