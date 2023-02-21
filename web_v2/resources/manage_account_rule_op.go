package resources

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func newManageAccountRuleOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageAccountRule
	switch details.Action {
	case xdr.ManageAccountRuleActionCreate:
		return &regources.ManageAccountRuleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_ACCOUNT_RULE),
			Attributes: &regources.ManageAccountRuleOpAttributes{
				Resource: details.CreateDetails.Resource,
				Action:   details.CreateDetails.Action,
				Forbids:  details.CreateDetails.IsForbid,
				Details:  details.CreateDetails.Details,
			},
			Relationships: &regources.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionUpdate:
		return &regources.ManageAccountRuleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_UPDATE_ACCOUNT_RULE),
			Attributes: &regources.ManageAccountRuleOpAttributes{
				Resource: details.UpdateDetails.Resource,
				Action:   details.UpdateDetails.Action,
				Forbids:  details.UpdateDetails.IsForbid,
				Details:  details.UpdateDetails.Details,
			},
			Relationships: &regources.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageAccountRuleActionRemove:
		return &regources.ManageAccountRuleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_REMOVE_ACCOUNT_RULE),
			Relationships: &regources.ManageAccountRuleOpRelationships{
				Rule: NewAccountRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
