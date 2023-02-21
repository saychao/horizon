package resources

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"gitlab.com/tokend/go/xdr"

	"github.com/saychao/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewFee - creates new instance of Fee from provided one.
func NewFee(record core2.Fee) regources.FeeRecord {
	hash := CalculateFeeHash(record)
	return regources.FeeRecord{
		Key: regources.Key{
			ID:   hash,
			Type: regources.FEES,
		},
		Attributes: regources.FeeRecordAttributes{
			Fixed:   regources.Amount(record.Fixed),
			Percent: regources.Amount(record.Percent),
			AppliedTo: regources.FeeAppliedTo{
				Asset:           record.Asset,
				Subtype:         record.Subtype,
				FeeType:         record.FeeType,
				FeeTypeExtended: xdr.FeeType(record.FeeType),
				LowerBound:      regources.Amount(record.LowerBound),
				UpperBound:      regources.Amount(record.UpperBound),
			},
		},
	}
}

// NewFeeKey - creates new Key for fee
func NewFeeKey(hash string) regources.Key {
	return regources.Key{
		ID:   hash,
		Type: regources.FEES,
	}
}

func NewCalculatedFeeKey(hash string) regources.Key {
	return regources.Key{
		ID:   hash,
		Type: regources.CALCULATED_FEE,
	}
}

func CalculateFeeHash(fee core2.Fee) string {
	lowerBound := regources.Amount(fee.LowerBound)
	upperBound := regources.Amount(fee.UpperBound)
	data := fmt.Sprintf("type:%d:asset:%s:subtype:%d:lower_bound:%s:upper_bound:%s",
		fee.FeeType,
		fee.Asset,
		fee.Subtype,
		lowerBound.String(),
		upperBound.String(),
	)

	if fee.AccountID != "" {
		data += fmt.Sprintf("account_id:%s", fee.AccountID)
	} else {
		data += fmt.Sprintf("account_role:%d", fee.AccountRole)
	}

	rawHash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(rawHash[:])
}
