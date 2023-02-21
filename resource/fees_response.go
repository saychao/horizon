package resource

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type FeesResponse struct {
	Fees map[xdr.AssetCode][]regources.FeeEntry `json:"fees"`
}
