package resources

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/spf13/cast"
	regources "gitlab.com/tokend/regources/generated"
)

// NewExternalSystemID creates new instance of ExternalSystemID from provided one
func NewExternalSystemID(extSysID core2.ExternalSystemID) *regources.ExternalSystemId {
	return &regources.ExternalSystemId{
		Key: regources.Key{
			ID:   cast.ToString(extSysID.ID),
			Type: regources.EXTERNAL_SYSTEM_IDS,
		},
		Attributes: regources.ExternalSystemIdAttributes{
			ExternalSystemType: extSysID.ExternalSystemType,
			Data: regources.ExternalSystemData{
				Data: extSysID.Data.Data,
				Type: extSysID.Data.Type,
			},
			IsDeleted: extSysID.IsDeleted,
			ExpiresAt: extSysID.ExpiresAt,
			BindedAt:  extSysID.BindedAt,
		},
	}
}
