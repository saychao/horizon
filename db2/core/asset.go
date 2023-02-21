package core

import (
	"encoding/json"
	"math"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const maximumTrailingDigits int32 = 6

type Asset struct {
	Code                 string `db:"code"`
	Policies             int32  `db:"policies"`
	Owner                string `db:"owner"`
	AvailableForIssuance uint64 `db:"available_for_issueance"`
	PreissuedAssetSigner string `db:"preissued_asset_signer"`
	MaxIssuanceAmount    uint64 `db:"max_issuance_amount"`
	Issued               uint64 `db:"issued"`
	LockedIssuance       uint64 `db:"locked_issuance"`
	PendingIssuance      uint64 `db:"pending_issuance"`
	Details              []byte `db:"details"`
	TrailingDigits       int32  `db:"trailing_digits"`
	Type                 uint64 `db:"type"`
}

func (a Asset) GetDetails() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(a.Details, &result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal asset details")
	}

	return result, nil
}

func (a Asset) GetMinimumAmount() int64 {
	nullDigits := maximumTrailingDigits - a.TrailingDigits
	if nullDigits < 0 {
		panic("Unexpected database state. Expected asset trailing digits be equal or less 6")
	}

	return int64(math.Pow10(int(nullDigits)))
}
