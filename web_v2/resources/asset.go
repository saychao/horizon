package resources

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

// NewAsset - creates new instance of Asset from provided one using history2 package.
func NewAssetV2(record history2.Asset) regources.Asset {
	return regources.Asset{
		Key: regources.Key{
			ID:   record.Code,
			Type: regources.ASSETS,
		},
		Attributes: regources.AssetAttributes{
			PreIssuanceAssetSigner: record.PreIssuanceAssetSigner,
			Details:                regources.Details(record.Details),
			MaxIssuanceAmount:      regources.Amount(record.MaxIssuanceAmount),
			AvailableForIssuance:   regources.Amount(record.AvailableForIssuance),
			Issued:                 regources.Amount(record.Issued),
			PendingIssuance:        regources.Amount(record.PendingIssuance),
			Policies:               xdr.AssetPolicy(record.Policies),
			TrailingDigits:         record.TrailingDigits,
			Type:                   record.Type,
			State:                  regources.AssetState(record.State),
		},
		Relationships: regources.AssetRelationships{
			Owner: NewAccountKey(record.Owner).AsRelation(),
		},
	}
}

// NewAsset - creates new instance of Asset from provided one.
func NewAsset(record core2.Asset) regources.Asset {
	return regources.Asset{
		Key: regources.Key{
			ID:   record.Code,
			Type: regources.ASSETS,
		},
		Attributes: regources.AssetAttributes{
			PreIssuanceAssetSigner: record.PreIssuanceAssetSigner,
			Details:                record.Details,
			MaxIssuanceAmount:      regources.Amount(record.MaxIssuanceAmount),
			AvailableForIssuance:   regources.Amount(record.AvailableForIssuance),
			Issued:                 regources.Amount(record.Issued),
			PendingIssuance:        regources.Amount(record.PendingIssuance),
			Policies:               xdr.AssetPolicy(record.Policies),
			TrailingDigits:         record.TrailingDigits,
			Type:                   record.Type,
			State:                  regources.AssetState(record.State),
		},
		Relationships: regources.AssetRelationships{
			Owner: NewAccountKey(record.Owner).AsRelation(),
		},
	}
}

// NewAssetKey - creates new Key for asset
func NewAssetKey(assetCode string) regources.Key {
	return regources.Key{
		ID:   assetCode,
		Type: regources.ASSETS,
	}
}
