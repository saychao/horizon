package reviewablerequest

import (
	"github.com/saychao/horizon/db2/history"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateDetails(request *history.ReviewableRequest, requestType xdr.ReviewableRequestType, h history.ReviewableRequestDetails) (
	d *regources.ReviewableRequestDetails, err error,
) {
	d = &regources.ReviewableRequestDetails{}
	d.RequestTypeName = requestType.ShortString()
	d.RequestType = int32(requestType)
	switch requestType {
	case xdr.ReviewableRequestTypeCreateAsset:
		d.AssetCreate, err = PopulateAssetCreationRequest(*h.AssetCreation)
		return
	case xdr.ReviewableRequestTypeUpdateAsset:
		d.AssetUpdate, err = PopulateAssetUpdateRequest(*h.AssetUpdate)
		return
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		d.PreIssuanceCreate, err = PopulatePreIssuanceRequest(*h.PreIssuanceCreate)
		return
	case xdr.ReviewableRequestTypeCreateIssuance:
		d.IssuanceCreate, err = PopulateIssuanceRequest(*h.IssuanceCreate)
		return
	case xdr.ReviewableRequestTypeCreateWithdraw:
		d.Withdraw, err = PopulateWithdrawalRequest(*h.Withdraw)
		return
	case xdr.ReviewableRequestTypeCreateSale:
		d.Sale, err = PopulateSaleCreationRequest(*h.Sale)
		return
	case xdr.ReviewableRequestTypeUpdateLimits:
		d.LimitsUpdate, err = PopulateLimitsUpdateRequest(*h.LimitsUpdate)
		return
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		d.AMLAlert, err = PopulateAmlAlertRequest(*h.AmlAlert)
		return
	case xdr.ReviewableRequestTypeChangeRole:
		d.ChangeRole, err = PopulateChangeRoleRequest(request, *h.ChangeRole)
		return
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		d.UpdateSaleDetails, err = PopulateUpdateSaleDetailsRequest(*h.UpdateSaleDetails)
		return
	case xdr.ReviewableRequestTypeCreateInvoice:
		d.Invoice, err = PopulateInvoiceRequest(*h.Invoice)
		return
	case xdr.ReviewableRequestTypeManageContract:
		d.Contract, err = PopulateContractRequest(*h.Contract)
		return
	case xdr.ReviewableRequestTypeCreateAtomicSwapAsk:
		d.AtomicSwapBidCreation, err = PopulateAtomicSwapAskCreationRequest(*h.AtomicSwapAskCreation)
		return
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		d.AtomicSwap, err = PopulateAtomicSwapRequest(*h.AtomicSwap)
		return
	case xdr.ReviewableRequestTypeCreatePoll:
		return
	case xdr.ReviewableRequestTypeKycRecovery:
		return
	case xdr.ReviewableRequestTypePerformRedemption:
		return
	case xdr.ReviewableRequestTypeDataCreation:
		return
	case xdr.ReviewableRequestTypeDataUpdate:
		return
	case xdr.ReviewableRequestTypeDataRemove:
		return
	case xdr.ReviewableRequestTypeCreateDeferredPayment:
		return
	case xdr.ReviewableRequestTypeCloseDeferredPayment:
		return
	default:
		return nil, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": requestType.String(),
		})
	}
}
