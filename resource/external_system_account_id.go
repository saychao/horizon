package resource

import (
	"time"

	"github.com/SafeRE-IT/horizon/db2/core"
	"gitlab.com/tokend/regources"
)

func PopulateExternalSystemAccountID(coreRecord core.ExternalSystemAccountID) (id regources.ExternalSystemAccountID) {
	switch coreRecord.ExternalSystemType {
	case 1:
		id.Type.Name = "Bitcoin"
		id.AssetCode = "BTC"
	case 2:
		id.Type.Name = "Ethereum"
		id.AssetCode = "ETH"
	}
	id.Type.Value = int32(coreRecord.ExternalSystemType)
	id.Data = coreRecord.Data

	// check out actions_account.go to find ExpiresAt default value
	if coreRecord.ExpiresAt != nil {
		expiresAt := time.Unix(*coreRecord.ExpiresAt, 0).Format(time.RFC3339)
		id.ExpiresAt = &expiresAt
	}
	return id
}
