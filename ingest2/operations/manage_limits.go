package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type manageLimitsOpHandler struct {
	effectsProvider
}

// Details returns details about manage limits operation
func (h *manageLimitsOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageLimitsOp := op.Body.MustManageLimitsOp()

	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeManageLimits,
		ManageLimits: &history2.ManageLimitsDetails{
			Action: manageLimitsOp.Details.Action,
		},
	}

	switch opDetails.ManageLimits.Action {
	case xdr.ManageLimitsActionCreate:
		creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

		opDetails.ManageLimits.Creation = &history2.ManageLimitsCreationDetails{
			AccountAddress:  creationDetails.AccountId.Address(), // Address() - smart, check for nil inside
			AccountRole:     creationDetails.AccountRole,
			StatsOpType:     creationDetails.StatsOpType,
			AssetCode:       string(creationDetails.AssetCode),
			IsConvertNeeded: creationDetails.IsConvertNeeded,
			DailyOut:        regources.Amount(creationDetails.DailyOut),
			WeeklyOut:       regources.Amount(creationDetails.WeeklyOut),
			MonthlyOut:      regources.Amount(creationDetails.MonthlyOut),
			AnnualOut:       regources.Amount(creationDetails.AnnualOut),
		}
	case xdr.ManageLimitsActionRemove:
		opDetails.ManageLimits.Removal = &history2.ManageLimitsRemovalDetails{
			LimitsID: int64(manageLimitsOp.Details.MustId()),
		}
	default:
		return history2.OperationDetails{}, errors.New("unexpected manage limits action")
	}

	return opDetails, nil
}

// ParticipantsEffects - returns source of the operation and account for which limits were managed if they are different
func (h *manageLimitsOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	_ xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	source := h.Participant(sourceAccountID)
	participants := []history2.ParticipantEffect{source}

	manageLimitsOp := opBody.MustManageLimitsOp()

	if manageLimitsOp.Details.Action != xdr.ManageLimitsActionCreate {
		return participants, nil
	}

	creationDetails := manageLimitsOp.Details.MustLimitsCreateDetails()

	if creationDetails.AccountId == nil {
		return participants, nil
	}

	limitsParticipant := h.Participant(*creationDetails.AccountId)

	if source.AccountID == limitsParticipant.AccountID {
		return participants, nil
	}

	participants = append(participants, limitsParticipant)

	return participants, nil
}
