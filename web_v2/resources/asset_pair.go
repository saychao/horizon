package resources

import (
	"fmt"

	"github.com/saychao/horizon/db2/core2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

// NewAssetPair - creates new instance of AssetPair from provided one.
func NewAssetPair(record core2.AssetPair) regources.AssetPair {
	return regources.AssetPair{
		Key: regources.Key{
			ID:   fmt.Sprintf("%s:%s", record.Base, record.Quote),
			Type: regources.ASSET_PAIRS,
		},
		Attributes: regources.AssetPairAttributes{
			Price:                   regources.Amount(record.CurrentPrice),
			Policies:                xdr.AssetPairPolicy(record.Policies),
			MaxPriceStep:            regources.Amount(record.MaxPriceStep),
			PhysicalPriceCorrection: regources.Amount(record.PhysicalPriceCorrection),
		},
	}
}
