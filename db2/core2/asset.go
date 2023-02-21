package core2

import (
	regources "gitlab.com/tokend/regources/generated"
	"math"
)

const maximumTrailingDigits uint32 = 6

// Asset - db representation of asset
type Asset struct {
	Code                   string            `db:"code"`
	Owner                  string            `db:"owner"`
	PreIssuanceAssetSigner string            `db:"preissued_asset_signer"`
	Details                regources.Details `db:"details"`
	MaxIssuanceAmount      uint64            `db:"max_issuance_amount"`
	AvailableForIssuance   uint64            `db:"available_for_issueance"`
	Issued                 uint64            `db:"issued"`
	PendingIssuance        uint64            `db:"pending_issuance"`
	Policies               int32             `db:"policies"`
	TrailingDigits         uint32            `db:"trailing_digits"`
	Type                   uint64            `db:"type"`
	State                  uint32            `db:"state"`
}

//GetMinimumAmount - returns min amount support for that asset
func (a Asset) GetMinimumAmount() int64 {
	nullDigits := maximumTrailingDigits - a.TrailingDigits
	if nullDigits < 0 {
		panic("Unexpected database state. Expected asset trailing digits be equal or less 6")
	}

	return int64(math.Pow10(int(nullDigits)))
}
