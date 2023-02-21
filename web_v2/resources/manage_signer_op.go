package resources

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func newManageSignerOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageSigner
	switch details.Action {
	case xdr.ManageSignerActionCreate:
		return &regources.ManageSignerOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_SIGNER),
			Attributes: &regources.ManageSignerOpAttributes{
				Details:  details.CreateDetails.Details,
				Weight:   details.CreateDetails.Weight,
				Identity: details.CreateDetails.Identity,
			},
			Relationships: &regources.ManageSignerOpRelationships{
				Role:   NewSignerRoleKey(details.CreateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionUpdate:
		return &regources.ManageSignerOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_UPDATE_SIGNER),
			Attributes: &regources.ManageSignerOpAttributes{
				Details:  details.UpdateDetails.Details,
				Weight:   details.UpdateDetails.Weight,
				Identity: details.UpdateDetails.Identity,
			},
			Relationships: &regources.ManageSignerOpRelationships{
				Role:   NewSignerRoleKey(details.UpdateDetails.RoleID).AsRelation(),
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	case xdr.ManageSignerActionRemove:
		return &regources.ManageSignerOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_REMOVE_SIGNER),
			Relationships: &regources.ManageSignerOpRelationships{
				Signer: NewSignerKey(details.PublicKey).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
