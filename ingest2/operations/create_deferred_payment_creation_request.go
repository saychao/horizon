package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createDeferredPaymentCreationRequestOpHandler struct {
	effectsProvider
}

// Details returns details about manage balance operation
func (h *createDeferredPaymentCreationRequestOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	oper := op.Body.MustCreateDeferredPaymentCreationRequestOp()
	details := history2.OperationDetails{
		Type: xdr.OperationTypeCreateDeferredPaymentCreationRequest,
		CreateDeferredPaymentCreationRequest: &history2.CreateDeferredPaymentCreationRequest{
			RequestID:          uint64(oper.RequestId),
			SourceBalance:      oper.Request.SourceBalance.AsString(),
			DestinationAccount: oper.Request.Destination.Address(),
			Amount:             regources.Amount(oper.Request.Amount),
			Details:            internal.MarshalCustomDetails(oper.Request.CreatorDetails),
			AllTasks:           (*uint32)(oper.AllTasks),
		},
	}

	return details, nil
}

// ParticipantsEffects returns `locked` effect
func (h *createDeferredPaymentCreationRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	oper := opBody.MustCreateDeferredPaymentCreationRequestOp().Request

	locked := h.effectsProvider.BalanceEffect(oper.SourceBalance,
		&history2.Effect{
			Type: history2.EffectTypeLocked,
			Locked: &history2.BalanceChangeEffect{
				Amount: regources.Amount(oper.Amount),
			},
		})

	return []history2.ParticipantEffect{locked}, nil
}
