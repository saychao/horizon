package history2

import (
	"encoding/json"
	regources "gitlab.com/tokend/regources/generated"
)

type Asset struct {
	Code                   string          `db:"code"`
	Owner                  string          `db:"owner"`
	PreIssuanceAssetSigner string          `db:"preissued_asset_signer"`
	Details                json.RawMessage `db:"details"`
	MaxIssuanceAmount      uint64          `db:"max_issuance_amount"`
	AvailableForIssuance   uint64          `db:"available_for_issuance"`
	Issued                 uint64          `db:"issued"`
	PendingIssuance        uint64          `db:"pending_issuance"`
	Policies               uint32          `db:"policies"`
	TrailingDigits         uint32          `db:"trailing_digits"`
	Type                   uint64          `db:"type"`

	State regources.AssetState `db:"state"`
}
