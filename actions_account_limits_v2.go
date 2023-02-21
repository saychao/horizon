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
// LimitsV2AccountShowAction: renders AccountTypeLimits for account
type LimitsV2AccountShowAction struct {
	Action
	AccountID string
	Result    resource.LimitsResponse
}

// JSON is a method for actions.JSON
func (action *LimitsV2AccountShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Result)
		},
	)
}
func (action *LimitsV2AccountShowAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *LimitsV2AccountShowAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *LimitsV2AccountShowAction) loadData() {
	account, err := action.CoreQ().Accounts().ByAddress(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load account by ID")
		action.Err = &problem.ServerError
		return
	}

	if account == nil {
		action.Err = &problem.NotFound
		return
	}

	limits, err := action.CoreQ().LimitsV2().ForAccount(account)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load limits for account")
		action.Err = &problem.ServerError
		return
	}

	stats, err := action.loadStatistics(account.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to load stats for account")
		action.Err = &problem.ServerError
		return
	}

	for _, limit := range limits {
		var result resource.LimitResponse
		result.Limit.Populate(limit)
		stat, ok := stats[statsKey{
			StatsOpType:     limit.StatsOpType,
			AssetCode:       limit.AssetCode,
			IsConvertNeeded: limit.IsConvertNeeded,
		}]

		if ok {
			result.Statistics.Populate(stat)
		} else {
			result.Statistics.IsConvertNeeded = limit.IsConvertNeeded
			result.Statistics.AssetCode = limit.AssetCode
			result.Statistics.StatsOpType = limit.StatsOpType
			result.Statistics.AccountId = action.AccountID
		}

		action.Result.Limits = append(action.Result.Limits, result)

	}
}

type statsKey struct {
	StatsOpType     int32
	AssetCode       string
	IsConvertNeeded bool
}

func (action *LimitsV2AccountShowAction) loadStatistics(accountID string) (map[statsKey]core.StatisticsV2Entry, error) {
	rawStats, err := action.CoreQ().StatisticsV2().ForAccount(accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load statistics for account")
	}

	result := map[statsKey]core.StatisticsV2Entry{}
	for i := range rawStats {
		result[statsKey{
			StatsOpType:     rawStats[i].StatsOpType,
			AssetCode:       rawStats[i].AssetCode,
			IsConvertNeeded: rawStats[i].IsConvertNeeded,
		}] = rawStats[i]
	}

	return result, nil
}
