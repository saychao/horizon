package reviewablerequest

import (
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateAssetUpdateRequest(histRequest history.AssetUpdateRequest) (
	*regources.AssetUpdateRequest, error,
) {
	return &regources.AssetUpdateRequest{
		Code:     histRequest.Asset,
		Policies: base.FlagFromXdrAssetPolicy(histRequest.Policies, xdr.AssetPolicyAll),
		Details:  histRequest.Details,
	}, nil
}
