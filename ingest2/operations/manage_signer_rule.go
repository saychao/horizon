package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageSignerRuleOpHandler struct {
	effectsProvider
}

// Details returns details about bind external system account operation
func (h *manageSignerRuleOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageSignerRuleOp := op.Body.MustManageSignerRuleOp()
	manageSignerRuleResult := opRes.MustManageSignerRuleResult().MustSuccess()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageSignerRule,
		ManageSignerRule: &history2.ManageSignerRuleDetails{
			RuleID: uint64(manageSignerRuleResult.RuleId),
			Action: manageSignerRuleOp.Data.Action,
		},
	}

	switch manageSignerRuleOp.Data.Action {
	case xdr.ManageSignerRuleActionCreate:
		details := manageSignerRuleOp.Data.MustCreateData()

		opDetails.ManageSignerRule.CreateDetails = &history2.CreateSignerRuleDetails{
			Details:    internal.MarshalCustomDetails(details.Details),
			IsReadOnly: details.IsReadOnly,
			IsForbid:   details.Forbids,
			IsDefault:  details.IsDefault,
			Resource:   details.Resource,
			Action:     details.Action,
		}
	case xdr.ManageSignerRuleActionUpdate:
		details := manageSignerRuleOp.Data.MustUpdateData()

		opDetails.ManageSignerRule.UpdateDetails = &history2.UpdateSignerRuleDetails{
			Details:   internal.MarshalCustomDetails(details.Details),
			IsForbid:  details.Forbids,
			IsDefault: details.IsDefault,
			Resource:  details.Resource,
			Action:    details.Action,
		}
	case xdr.ManageSignerRuleActionRemove:
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account role")
	}

	return opDetails, nil
}

// ParticipantsEffects returns only source without effects
func (h *manageSignerRuleOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	return []history2.ParticipantEffect{h.Participant(sourceAccountID)}, nil
}
