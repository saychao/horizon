package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageSignerOpHandler struct {
	effectsProvider
}

// Details returns details about manage signer operation
func (h *manageSignerOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageSignerOp := op.Body.MustManageSignerOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageSigner,
		ManageSigner: &history2.ManageSignerDetails{
			Action: manageSignerOp.Data.Action,
		},
	}

	var publicKey xdr.AccountId

	switch manageSignerOp.Data.Action {
	case xdr.ManageSignerActionCreate:
		details := manageSignerOp.Data.MustCreateData()
		publicKey = xdr.AccountId(details.PublicKey)

		opDetails.ManageSigner.CreateDetails = &history2.UpdateSignerDetails{
			Details:  internal.MarshalCustomDetails(xdr.Longstring(details.Details)),
			RoleID:   uint64(details.RoleId),
			Identity: uint32(details.Identity),
			Weight:   uint32(details.Weight),
		}
	case xdr.ManageSignerActionUpdate:
		details := manageSignerOp.Data.MustUpdateData()
		publicKey = xdr.AccountId(details.PublicKey)

		opDetails.ManageSigner.UpdateDetails = &history2.UpdateSignerDetails{
			Details:  internal.MarshalCustomDetails(xdr.Longstring(details.Details)),
			RoleID:   uint64(details.RoleId),
			Identity: uint32(details.Identity),
			Weight:   uint32(details.Weight),
		}
	case xdr.ManageSignerActionRemove:
		publicKey = xdr.AccountId(manageSignerOp.Data.RemoveData.PublicKey)
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account role")
	}

	opDetails.ManageSigner.PublicKey = publicKey.Address()

	return opDetails, nil
}

// ParticipantsEffects returns only source without effects
func (h *manageSignerOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
