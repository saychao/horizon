package resource

import (
	"strconv"

	"github.com/SafeRE-IT/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type Asset struct {
	Code                 string `json:"code"`
	Type                 uint64 `json:"type"`
	Owner                string `json:"owner"`
	AvailableForIssuance string `json:"available_for_issuance"`
	PreissuedAssetSigner string `json:"preissued_asset_signer"`
	MaxIssuanceAmount    string `json:"max_issuance_amount"`
	Issued               string `json:"issued"`
	PendingIssuance      string `json:"pending_issuance"`
	TrailingDigitsCount  string `json:"trailing_digits_count"`
	Policies
	Details map[string]interface{} `json:"details"`
	Sales   []Sale                 `json:"sales,omitempty"`
}

func (a *Asset) Populate(asset *core.Asset) {
	a.Code = asset.Code
	a.Owner = asset.Owner
	a.PreissuedAssetSigner = asset.PreissuedAssetSigner

	a.AvailableForIssuance = amount.StringU(asset.AvailableForIssuance)
	a.MaxIssuanceAmount = amount.StringU(asset.MaxIssuanceAmount)
	a.PendingIssuance = amount.StringU(asset.PendingIssuance)
	a.Issued = amount.StringU(asset.Issued)
	a.TrailingDigitsCount = strconv.Itoa(int(asset.TrailingDigits))
	a.Type = asset.Type

	a.Policies.Populate(*asset)
	a.Details, _ = asset.GetDetails()
}

func PopulateAssetPair(asset core.AssetPair) regources.AssetPair {
	return regources.AssetPair{
		Base:                    asset.BaseAsset,
		Quote:                   asset.QuoteAsset,
		CurrentPrice:            regources.Amount(asset.CurrentPrice),
		PhysicalPrice:           regources.Amount(asset.PhysicalPrice),
		PhysicalPriceCorrection: regources.Amount(asset.PhysicalPriceCorrection),
		MaxPriceStep:            regources.Amount(asset.MaxPriceStep),
		Policy:                  asset.Policies,
		Policies:                PopulatePolicies(asset.Policies),
	}
}

func PopulatePolicies(policy int32) []regources.Policy {
	result := make([]regources.Policy, 0)

	for _, p := range xdr.AssetPairPolicyAll {
		if (int32(p) & policy) != 0 {
			result = append(result, regources.Policy{
				Name:  p.String(),
				Value: int32(p),
			})
		}
	}

	return result
}
