package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createRedemptionRequestOpHandler struct {
	effectsProvider
}

func (h *createRedemptionRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createRedemptionRequestOp := op.Body.MustCreateRedemptionRequestOp()
	response := opRes.MustCreateRedemptionRequestResult().MustRedemptionResponse()

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateRedemptionRequest,
		Redemption: &history2.RedemptionDetails{
			BalanceFrom: createRedemptionRequestOp.RedemptionRequest.SourceBalanceId.AsString(),
			AccountTo:   createRedemptionRequestOp.RedemptionRequest.Destination.Address(),
			Asset:       string(response.Asset),
			Amount:      regources.Amount(createRedemptionRequestOp.RedemptionRequest.Amount),
			Details:     regources.Details(createRedemptionRequestOp.RedemptionRequest.CreatorDetails),
			RequestDetails: history2.RequestDetails{
				IsFulfilled: response.Fulfilled,
				RequestID:   int64(response.RequestId),
			},
		},
	}, nil
}

func (h *createRedemptionRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	op := opBody.MustCreateRedemptionRequestOp()

	return []history2.ParticipantEffect{
		h.BalanceEffect(op.RedemptionRequest.SourceBalanceId, &history2.Effect{
			Type: history2.EffectTypeLocked,
			Locked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(op.RedemptionRequest.Amount),
			}}),
	}, nil
}
