package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createCloseDeferredPaymentRequestOpHandler struct {
	effectsProvider
	defPaymentProvider
}

// Details returns details about manage balance operation
func (h *createCloseDeferredPaymentRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateCloseDeferredPaymentRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCreateCloseDeferredPaymentRequest,
		CreateCloseDeferredPaymentRequest: &history2.CreateCloseDeferredPaymentRequest{
			RequestID: uint64(oper.RequestId),
			Amount:    regources.Amount(oper.Request.Amount),
			Details:   internal.MarshalCustomDetails(oper.Request.CreatorDetails),
			AllTasks:  (*uint32)(oper.AllTasks),
		},
	}

	switch oper.Request.Destination.Type {
	case xdr.CloseDeferredPaymentDestinationTypeAccount:
		details.CreateCloseDeferredPaymentRequest.DestinationAccount = oper.Request.Destination.AccountId.Address()
	case xdr.CloseDeferredPaymentDestinationTypeBalance:
		details.CreateCloseDeferredPaymentRequest.DestinationBalance = oper.Request.Destination.BalanceId.AsString()
	}

	return details, nil
}

// ParticipantsEffects returns `funded` and `charged` effects
func (h *createCloseDeferredPaymentRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	requestResult := opRes.MustCreateCloseDeferredPaymentRequestResult().MustSuccess()

	if !requestResult.Fulfilled {
		return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
	}

	closeRequest := opBody.MustCreateCloseDeferredPaymentRequestOp().Request

	dp := h.MustDeferredPayment(int64(requestResult.DeferredPaymentId))
	var sb xdr.BalanceId
	sb.SetString(dp.SourceBalance)

	// just unlock funds if deferred payment has been returned to the source balance
	if requestResult.ExtendedResult.DestinationBalance.Equals(sb) {
		unlocked := h.effectsProvider.BalanceEffect(sb,
			&history2.Effect{
				Type: history2.EffectTypeUnlocked,
				Unlocked: &history2.BalanceChangeEffect{
					Amount: regources.Amount(closeRequest.Amount),
				},
			})

		return []history2.ParticipantEffect{unlocked}, nil
	}

	funded := h.effectsProvider.BalanceEffect(requestResult.ExtendedResult.DestinationBalance,
		&history2.Effect{
			Type: history2.EffectTypeFunded,
			Funded: &history2.BalanceChangeEffect{
				Amount: regources.Amount(closeRequest.Amount),
			},
		})

	charged := h.effectsProvider.BalanceEffect(sb,
		&history2.Effect{
			Type: history2.EffectTypeChargedFromLocked,
			ChargedFromLocked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(closeRequest.Amount),
			},
		})

	return []history2.ParticipantEffect{funded, charged}, nil
}
