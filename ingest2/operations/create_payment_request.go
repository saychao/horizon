package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"github.com/saychao/horizon/utf8"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createPaymentRequestOpHandler struct {
	paymentHandler *paymentOpHandler
}

// Details returns details about create limits request operation
func (h *createPaymentRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	createPaymentRequestOp := op.Body.MustCreatePaymentRequestOp()
	createPaymentRequestRes := opRes.MustCreatePaymentRequestResult().MustSuccess()
	paymentOp := createPaymentRequestOp.Request.PaymentOp
	request := createPaymentRequestOp.Request
	var allTasks *uint32
	if createPaymentRequestOp.AllTasks != nil {
		tasks := uint32(*createPaymentRequestOp.AllTasks)
		allTasks = &tasks
	}

	creatorDetails := regources.Details("{}")
	switch request.Ext.V {
	case xdr.LedgerVersionMovementRequestsDetails:
		creatorDetails = internal.MarshalCustomDetails(request.Ext.MustCreatorDetails())
	case xdr.LedgerVersionEmptyVersion:
	default:
		panic(errors.From(errors.New("unexpected version of payment request"), logan.F{
			"ledger_version": request.Ext.V,
		}))
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreatePaymentRequest,
		CreatePaymentRequest: &history2.CreatePaymentRequestDetails{
			PaymentDetails: history2.PaymentRequestDetails{
				CreatorDetails:          creatorDetails,
				AccountFrom:             op.Source.Address(),
				BalanceFrom:             paymentOp.SourceBalanceId.AsString(),
				Amount:                  regources.Amount(paymentOp.Amount),
				SourceFee:               internal.FeeFromXdr(paymentOp.FeeData.SourceFee),
				DestinationFee:          internal.FeeFromXdr(paymentOp.FeeData.DestinationFee),
				SourcePayForDestination: paymentOp.FeeData.SourcePaysForDest,
				Subject:                 string(paymentOp.Subject),
				Reference:               utf8.Scrub(string(paymentOp.Reference)),
			},
			AllTasks: allTasks,
			RequestDetails: history2.RequestDetails{
				RequestID:   int64(createPaymentRequestRes.RequestId),
				IsFulfilled: createPaymentRequestRes.Fulfilled,
			},
		},
	}, nil
}

func (h *createPaymentRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	paymentOp := opBody.MustCreatePaymentRequestOp().Request.PaymentOp
	result := opRes.MustCreatePaymentRequestResult().MustSuccess()

	if !result.Fulfilled {
		return []history2.ParticipantEffect{h.paymentHandler.Participant(sourceAccountID)}, nil
	}

	if result.PaymentResult == nil {
		return nil, errors.New("unexpected nil payment result")
	}

	return h.paymentHandler.participantEffects(paymentOp,
		result.PaymentResult.MustPaymentResponse(), sourceAccountID)
}
