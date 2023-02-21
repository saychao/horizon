package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type manageSaleHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// Details returns details about payout operation
func (h *manageSaleHandler) Details(op rawOperation, res xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageSale := op.Body.MustManageSaleOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageSale,
		ManageSale: &history2.ManageSaleDetails{
			SaleId: uint64(manageSale.SaleId),
			Action: manageSale.Data.Action,
		},
	}, nil
}

// ParticipantsEffects returns `charged` and `funded` effects
func (h *manageSaleHandler) ParticipantsEffects(opBody xdr.OperationBody,
	res xdr.OperationResultTr, sourceAccountID xdr.AccountId, changes []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	manageSale := opBody.MustManageSaleOp().Data
	if manageSale.Action != xdr.ManageSaleActionCancel {
		return []history2.ParticipantEffect{h.manageOfferOpHandler.Participant(sourceAccountID)}, nil
	}

	return h.manageOfferOpHandler.getDeletedOffersEffect(changes), nil
}
