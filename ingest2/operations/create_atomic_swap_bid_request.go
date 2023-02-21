package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createAtomicSwapBidRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create atomic swap request operation
func (h *createAtomicSwapBidRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	aSwapRequest := op.Body.MustCreateAtomicSwapBidRequestOp().Request
	successRes := opRes.MustCreateAtomicSwapBidRequestResult().MustSuccess()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAtomicSwapBidRequest,
		CreateAtomicSwapBidRequest: &history2.CreateAtomicSwapBidRequestDetails{
			AskID:      int64(aSwapRequest.AskId),
			BaseAmount: regources.Amount(aSwapRequest.BaseAmount),
			QuoteAsset: string(aSwapRequest.QuoteAsset),
			RequestDetails: history2.RequestDetails{
				RequestID: int64(successRes.RequestId),
			},
		},
	}, nil
}
