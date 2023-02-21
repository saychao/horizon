package internal

import (
	"time"

	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

//FeeFromXdr - converts xdr fee to regources fee
func FeeFromXdr(data xdr.Fee) regources.Fee {
	return regources.Fee{
		Fixed:             regources.Amount(data.Fixed),
		CalculatedPercent: regources.Amount(data.Percent),
	}
}

//TimeFromXdr - converts xdr.Uint64 to unix utc timestamp
func TimeFromXdr(data xdr.Uint64) time.Time {
	return time.Unix(int64(data), 0).UTC()
}
