package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type manageAccountSpecificRuleOpHandler struct {
	manageOfferOpHandler *manageOfferOpHandler
}

// Details returns details about bind external system account operation
func (h *manageAccountSpecificRuleOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageRuleOp := op.Body.MustManageAccountSpecificRuleOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageAccountSpecificRule,
		ManageAccountSpecificRule: &history2.ManageAccountSpecificRuleDetails{
			Action: manageRuleOp.Data.Action,
		},
	}

	switch manageRuleOp.Data.Action {
	case xdr.ManageAccountSpecificRuleActionCreate:
		opDetails.ManageAccountSpecificRule.RuleID =
			uint64(opRes.MustManageAccountSpecificRuleResult().MustSuccess().RuleId)
		details := manageRuleOp.Data.MustCreateData()

		creationDetails := &history2.CreateAccountSpecificRuleDetails{
			LedgerKey: details.LedgerKey,
			Forbids:   details.Forbids,
			AccountID: details.AccountId.Address(), // Address() smart
		}

		opDetails.ManageAccountSpecificRule.CreateDetails = creationDetails
	case xdr.ManageAccountSpecificRuleActionRemove:
		opDetails.ManageAccountSpecificRule.RuleID = uint64(manageRuleOp.Data.MustRemoveData().RuleId)
	default:
		return history2.OperationDetails{}, errors.New("Unexpected action on manage account specific rule")
	}

	return opDetails, nil
}

func (h *manageAccountSpecificRuleOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, ledgerChanges []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	opData := opBody.MustManageAccountSpecificRuleOp().Data

	switch opData.Action {
	case xdr.ManageAccountSpecificRuleActionCreate:
		return []history2.ParticipantEffect{h.manageOfferOpHandler.Participant(sourceAccountID)}, nil
	case xdr.ManageAccountSpecificRuleActionRemove:
		return append(h.manageOfferOpHandler.getDeletedOffersEffect(ledgerChanges),
			h.manageOfferOpHandler.Participant(sourceAccountID)), nil
	default:
		return nil, errors.From(errors.New("unexpected manage account specific rule action"),
			map[string]interface{}{
				"action_i": int32(opData.Action),
			})
	}
}
