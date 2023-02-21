package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type cancelDeferredPaymentCreationRequestOpHandler struct {
	effectsProvider
	reviewableRequestsStorage
}

// Details returns details about manage balance operation
func (h *cancelDeferredPaymentCreationRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	cancelDeferredPaymentCreationCreationRequestOp := op.Body.MustCancelDeferredPaymentCreationRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCancelDeferredPaymentCreationRequest,
		CancelDeferredPaymentCreationRequest: &history2.CancelDeferredPaymentCreationRequest{
			RequestID: uint64(cancelDeferredPaymentCreationCreationRequestOp.RequestId),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `unlocked` effect
func (h *cancelDeferredPaymentCreationRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	cancelRequestResult := opRes.MustCancelDeferredPaymentCreationRequestResult()

	if cancelRequestResult.Code != xdr.CancelDeferredPaymentCreationRequestResultCodeSuccess {
		return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
	}

	createRequestID := uint64(opBody.MustCancelDeferredPaymentCreationRequestOp().RequestId)

	createRequest, err := h.GetByID(createRequestID)
	if err != nil {
		return []history2.ParticipantEffect{}, err
	}

	if createRequest == nil {
		return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
	}

	var sb xdr.BalanceId
	sb.SetString(createRequest.Details.CreateDeferredPayment.SourceBalance)

	unlocked := h.effectsProvider.BalanceEffect(sb,
		&history2.Effect{
			Type: history2.EffectTypeUnlocked,
			Unlocked: &history2.BalanceChangeEffect{
				Amount: createRequest.Details.CreateDeferredPayment.Amount,
			},
		})

	return []history2.ParticipantEffect{unlocked}, nil
}
