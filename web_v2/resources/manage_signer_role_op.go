package resources

import (
	"errors"

	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func newManageSignerRoleOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageSignerRole
	switch details.Action {
	case xdr.ManageSignerRoleActionCreate:
		return &regources.ManageSignerRoleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_SIGNER_ROLE),
			Attributes: &regources.ManageSignerRoleOpAttributes{
				Details:    details.CreateDetails.Details,
				IsReadOnly: details.CreateDetails.IsReadOnly,
			},
			Relationships: regources.ManageSignerRoleOpRelationships{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, regources.SIGNER_RULES),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionUpdate:
		return &regources.ManageSignerRoleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_UPDATE_SIGNER_ROLE),
			Attributes: &regources.ManageSignerRoleOpAttributes{
				Details: details.UpdateDetails.Details,
			},
			Relationships: regources.ManageSignerRoleOpRelationships{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, regources.SIGNER_RULES),
				Role:  NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRoleActionRemove:
		return &regources.ManageSignerRoleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_REMOVE_SIGNER_ROLE),
			Relationships: regources.ManageSignerRoleOpRelationships{
				Role: NewSignerRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
