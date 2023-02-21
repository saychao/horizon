package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageAccountRuleOpHandler struct {
	effectsProvider
}

// Details returns details about bind external system account operation
func (h *manageAccountRuleOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageAccountRuleOp := op.Body.MustManageAccountRuleOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAccountRule,
		ManageAccountRule: &history2.ManageAccountRuleDetails{
			Action: manageAccountRuleOp.Data.Action,
		},
	}

	switch manageAccountRuleOp.Data.Action {
	case xdr.ManageAccountRuleActionCreate:
		opDetails.ManageAccountRule.RuleID =
			uint64(opRes.MustManageAccountRuleResult().MustSuccess().RuleId)
		details := manageAccountRuleOp.Data.MustCreateData()

		creationDetails := &history2.UpdateAccountRuleDetails{
			Resource: details.Resource,
			Action:   details.Action,
			IsForbid: details.Forbids,
			Details:  internal.MarshalCustomDetails(details.Details),
		}

		opDetails.ManageAccountRule.CreateDetails = creationDetails
	case xdr.ManageAccountRuleActionUpdate:
		details := manageAccountRuleOp.Data.MustUpdateData()
		opDetails.ManageAccountRule.RuleID = uint64(details.RuleId)

		updateDetails := &history2.UpdateAccountRuleDetails{
			Resource: details.Resource,
			Action:   details.Action,
			IsForbid: details.Forbids,
			Details:  internal.MarshalCustomDetails(details.Details),
		}

		opDetails.ManageAccountRule.UpdateDetails = updateDetails
	case xdr.ManageAccountRuleActionRemove:
		opDetails.ManageAccountRule.RuleID = uint64(manageAccountRuleOp.Data.MustRemoveData().RuleId)
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account rule")
	}

	return opDetails, nil
}

// ParticipantsEffects returns only source without effects
func (h *manageAccountRuleOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
