package resources

import (
	"time"

	"github.com/saychao/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewLicenseInfo - creates new instance of licenseInfo
func NewLicenseInfo(record core2.License, adminCount int64) regources.LicenseInfo {

	licenseInfo := regources.LicenseInfo{
		Key: regources.NewKeyInt64(record.ID, regources.LICENSE_INFO),
		Attributes: regources.LicenseInfoAttributes{
			CurrentAdminCount: adminCount,
			AdminCount:        record.AdminCount,
			DueDate:           time.Unix(record.DueDate, 0).UTC(),
		},
	}
	return licenseInfo
}
