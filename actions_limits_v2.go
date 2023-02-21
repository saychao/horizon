package horizon

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// This file contains the actions:
//
// AccountTypeLimitsShowAction: renders AccountTypeLimits for operationType
type LimitsV2ShowAction struct {
	Action
	StatsOpType int32
	Asset       string
	AccountID   *string
	AccountType *int32
	Result      resource.LimitsV2
}

// JSON is a method for actions.JSON
func (action *LimitsV2ShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Result)
		},
	)
}
func (action *LimitsV2ShowAction) loadParams() {
	action.StatsOpType = action.GetInt32("stats_op_type")
	action.Asset = action.GetNonEmptyString("asset")
	accountID := action.GetString("account_id")
	if accountID != "" {
		action.AccountID = &accountID
	}
	accountType := action.GetInt32("account_type")
	if accountType != 0 {
		action.AccountType = &accountType
	}

	if (action.AccountID != nil) && (action.AccountType != nil) {
		action.SetInvalidField("account_id", errors.New("It's not allowed to specify both account_id and account type"))
		return
	}
}

func (action *LimitsV2ShowAction) loadData() {
	q := action.CoreQ().LimitsV2().ForStatsOpType(action.StatsOpType).ForAsset(action.Asset)
	if action.AccountID != nil {
		q = q.ForAccountID(*action.AccountID)
	} else if action.AccountType != nil {
		q = q.ForAccountType(*action.AccountType)
	} else {
		q = q.Global()
	}

	result, err := q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to select limits")
		action.Err = &problem.ServerError
		return
	}

	if len(result) > 1 {
		action.Log.WithField("limits", result).Error("Expected 1 limit at max")
		action.Err = &problem.ServerError
		return
	}

	if len(result) == 0 {
		result = append(result, core.DefaultLimits(action.AccountType, action.AccountID, action.StatsOpType, action.Asset))
	}

	action.Result.Populate(result[0])
}
