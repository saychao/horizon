package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageSignerRoleOpHandler struct {
	effectsProvider
}

// Details returns details about bind external system account operation
func (h *manageSignerRoleOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageSignerRoleOp := op.Body.MustManageSignerRoleOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageSignerRole,
		ManageSignerRole: &history2.ManageSignerRoleDetails{
			Action: manageSignerRoleOp.Data.Action,
		},
	}

	manageSignerRoleResult := opRes.MustManageSignerRoleResult().MustSuccess()

	switch manageSignerRoleOp.Data.Action {
	case xdr.ManageSignerRoleActionCreate:
		opDetails.ManageSignerRole.RoleID =
			uint64(manageSignerRoleResult.RoleId)
		details := manageSignerRoleOp.Data.MustCreateData()

		creationDetails := &history2.CreateSignerRoleDetails{
			Details:    internal.MarshalCustomDetails(details.Details),
			IsReadOnly: details.IsReadOnly,
		}

		for _, id := range details.RuleIDs {
			creationDetails.RuleIDs = append(creationDetails.RuleIDs, uint64(id))
		}

		opDetails.ManageSignerRole.CreateDetails = creationDetails
	case xdr.ManageSignerRoleActionUpdate:
		details := manageSignerRoleOp.Data.MustUpdateData()
		opDetails.ManageSignerRole.RoleID = uint64(details.RoleId)

		updateDetails := &history2.UpdateSignerRoleDetails{
			Details: internal.MarshalCustomDetails(details.Details),
		}

		for _, id := range details.RuleIDs {
			updateDetails.RuleIDs = append(updateDetails.RuleIDs, uint64(id))
		}

		opDetails.ManageSignerRole.UpdateDetails = updateDetails
	case xdr.ManageSignerRoleActionRemove:
		opDetails.ManageSignerRole.RoleID = uint64(manageSignerRoleOp.Data.MustRemoveData().RoleId)
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account role")
	}

	return opDetails, nil
}

// ParticipantsEffects returns only source without effects
func (h *manageSignerRoleOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
