package resources

import (
	"errors"

	"strconv"

	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func newManageAccountRoleOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageAccountRole
	switch details.Action {
	case xdr.ManageAccountRoleActionCreate:
		return &regources.ManageAccountRoleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ACCOUNT_ROLE),
			Attributes: &regources.ManageAccountRoleOpAttributes{
				Details: details.CreateDetails.Details,
			},
			Relationships: regources.ManageAccountRoleOpRelationships{
				Rules: idsAsRelations(details.CreateDetails.RuleIDs, regources.ACCOUNT_RULES),
				Role:  NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRoleActionUpdate:
		return &regources.ManageAccountRoleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_UPDATE_ACCOUNT_ROLE),
			Attributes: &regources.ManageAccountRoleOpAttributes{
				Details: details.UpdateDetails.Details,
			},
			Relationships: regources.ManageAccountRoleOpRelationships{
				Rules: idsAsRelations(details.UpdateDetails.RuleIDs, regources.ACCOUNT_RULES),
				Role:  NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRoleActionRemove:
		return &regources.ManageAccountRoleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_REMOVE_ACCOUNT_ROLE),
			Relationships: regources.ManageAccountRoleOpRelationships{
				Role: NewAccountRoleKey(details.RoleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}

func idsAsRelations(ids []uint64, resourceType regources.ResourceType) *regources.RelationCollection {
	keys := make([]regources.Key, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, regources.Key{
			ID:   strconv.FormatUint(id, 10),
			Type: resourceType,
		})
	}

	return &regources.RelationCollection{
		Data: keys,
	}
}
