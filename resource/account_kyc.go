package resource

import (
	"encoding/json"

	"github.com/saychao/horizon/db2/core"
)

type AccountKYC struct {
	KYCData map[string]interface{} `db:"account_kyc_data"`
}

// Populate fills out the fields of the account KYC
func (ak *AccountKYC) Populate(row core.AccountKYC) {
	var kycData map[string]interface{}
	_ = json.Unmarshal([]byte(row.KYCData.String), &kycData)
	ak.KYCData = kycData
}
