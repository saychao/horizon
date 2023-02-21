package resources

import (
	"fmt"

	regources "gitlab.com/tokend/regources/generated"
)

func NewSaleWhitelistKey(ruleID uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprintf("%d", ruleID),
		Type: regources.SALE_WHITELIST,
	}
}

func NewSaleWhitelist(ruleId uint64, address string) regources.SaleWhitelist {
	return regources.SaleWhitelist{
		Key: NewSaleWhitelistKey(ruleId),
		Relationships: regources.SaleWhitelistRelationships{
			Participant: NewAccountKey(address).AsRelation(),
		},
	}
}
