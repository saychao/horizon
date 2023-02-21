package resources

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func newManageSignerRuleOp(op history2.Operation) regources.Resource {
	details := op.Details.ManageSignerRule
	switch details.Action {
	case xdr.ManageSignerRuleActionCreate:
		return &regources.ManageSignerRuleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_CREATE_SIGNER_RULE),
			Attributes: &regources.ManageSignerRuleOpAttributes{
				Resource:   details.CreateDetails.Resource,
				Action:     details.CreateDetails.Action,
				Forbids:    details.CreateDetails.IsForbid,
				IsDefault:  details.CreateDetails.IsDefault,
				IsReadOnly: details.CreateDetails.IsReadOnly,
				Details:    details.CreateDetails.Details,
			},
			Relationships: &regources.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionUpdate:
		return &regources.ManageSignerRuleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_UPDATE_SIGNER_RULE),
			Attributes: &regources.ManageSignerRuleOpAttributes{
				Resource:  details.UpdateDetails.Resource,
				Action:    details.UpdateDetails.Action,
				Forbids:   details.UpdateDetails.IsForbid,
				IsDefault: details.UpdateDetails.IsDefault,
				Details:   details.UpdateDetails.Details,
			},
			Relationships: &regources.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	case xdr.ManageSignerRuleActionRemove:
		return &regources.ManageSignerRuleOp{
			Key: regources.NewKeyInt64(op.ID, regources.OPERATIONS_REMOVE_SIGNER_RULE),
			Relationships: &regources.ManageSignerRuleOpRelationships{
				Rule: NewSignerRoleKey(details.RuleID).AsRelation(),
			},
		}
	default:
		panic(errors.New("unexpected manage account role action"))
	}
}
