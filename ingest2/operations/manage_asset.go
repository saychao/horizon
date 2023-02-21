package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type manageAssetOpHandler struct {
	effectsProvider
}

// Details returns details about manage asset operation
func (h *manageAssetOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAssetOp := op.Body.MustManageAssetOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAsset,
		ManageAsset: &history2.ManageAssetDetails{
			RequestID: int64(manageAssetOp.RequestId),
			Action:    manageAssetOp.Request.Action,
		},
	}

	if manageAssetOp.RequestId == 0 {
		opDetails.ManageAsset.RequestID = int64(opRes.MustManageAssetResult().MustSuccess().RequestId)
	}

	switch opDetails.ManageAsset.Action {
	case xdr.ManageAssetActionCreateAssetCreationRequest:
		creationDetails := manageAssetOp.Request.MustCreateAssetCreationRequest().CreateAsset

		policies := xdr.AssetPolicy(creationDetails.Policies)

		opDetails.ManageAsset.AssetCode = string(creationDetails.Code)
		opDetails.ManageAsset.Type = uint64(creationDetails.Type)
		opDetails.ManageAsset.CreatorDetails = internal.MarshalCustomDetails(creationDetails.CreatorDetails)
		opDetails.ManageAsset.Policies = &policies
		opDetails.ManageAsset.PreIssuanceSigner = creationDetails.PreissuedAssetSigner.Address()
		opDetails.ManageAsset.MaxIssuanceAmount = regources.Amount(creationDetails.MaxIssuanceAmount)
	case xdr.ManageAssetActionCreateAssetUpdateRequest:
		updateDetails := manageAssetOp.Request.MustCreateAssetUpdateRequest().UpdateAsset

		policies := xdr.AssetPolicy(updateDetails.Policies)

		opDetails.ManageAsset.AssetCode = string(updateDetails.Code)
		opDetails.ManageAsset.CreatorDetails = internal.MarshalCustomDetails(updateDetails.CreatorDetails)
		opDetails.ManageAsset.Policies = &policies
	case xdr.ManageAssetActionCancelAssetRequest:
	case xdr.ManageAssetActionChangePreissuedAssetSigner:
		data := manageAssetOp.Request.MustChangePreissuedSigner()

		opDetails.ManageAsset.AssetCode = string(data.Code)
		opDetails.ManageAsset.PreIssuanceSigner = data.AccountId.Address()
	case xdr.ManageAssetActionUpdateMaxIssuance:
		data := manageAssetOp.Request.MustUpdateMaxIssuance()

		opDetails.ManageAsset.AssetCode = string(data.AssetCode)
		opDetails.ManageAsset.MaxIssuanceAmount = regources.Amount(data.MaxIssuanceAmount)
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage asset action")
	}

	return opDetails, nil
}
