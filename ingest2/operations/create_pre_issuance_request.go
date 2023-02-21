package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createPreIssuanceRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create pre issuance request operation
func (h *createPreIssuanceRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	preissuanceRequest := op.Body.MustCreatePreIssuanceRequest().Request
	successResult := opRes.MustCreatePreIssuanceRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreatePreissuanceRequest,
		CreatePreIssuanceRequest: &history2.CreatePreIssuanceRequestDetails{
			AssetCode:   string(preissuanceRequest.Asset),
			Amount:      regources.Amount(preissuanceRequest.Amount),
			RequestID:   int64(successResult.RequestId),
			IsFulfilled: successResult.Fulfilled,
		},
	}, nil
}
