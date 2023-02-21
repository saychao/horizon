package resource

import (
	"encoding/json"

	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

type AssetEntry struct {
	Code     string                 `json:"code"`
	Owner    string                 `json:"owner"`
	Details  map[string]interface{} `json:"details"`
	Policies Policies               `json:"policies"`

	PreissuedAssetSigner  string `json:"preissued_asset_signer"`
	AvailableForIssueance string `json:"available_for_issueance"`
	Issued                string `json:"issued"`
	MaxIssuanceAmount     string `json:"max_issuance_amount"`
}

func (r *AssetEntry) Populate(entry xdr.AssetEntry) {
	r.Code = string(entry.Code)
	r.Owner = entry.Owner.Address()

	r.Details = make(map[string]interface{})
	//entry.Details is user data and we doesn't care about if it's not valid
	json.Unmarshal([]byte(entry.Details), &r.Details)
	r.Policies.PopulateFromInt32(int32(entry.Policies))
	r.PreissuedAssetSigner = entry.PreissuedAssetSigner.Address()
	r.AvailableForIssueance = amount.String(int64(entry.AvailableForIssueance))
	r.Issued = amount.String(int64(entry.Issued))
	r.MaxIssuanceAmount = amount.String(int64(entry.MaxIssuanceAmount))
}

type LedgerKeyAsset struct {
	AssetCode string `json:"asset_code"`
}

func (r *LedgerKeyAsset) Populate(entry xdr.LedgerKeyAsset) {
	r.AssetCode = string(entry.Code)
}
