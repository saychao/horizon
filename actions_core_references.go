package horizon

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
)

type CoreReferencesAction struct {
	Action

	accountID string
	reference string

	records []core.Reference
}

func (action *CoreReferencesAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		func() {
			response := map[string]interface{}{
				"data": action.records,
			}
			hal.Render(action.W, response)
		},
	)
}

func (action *CoreReferencesAction) loadParams() {
	action.accountID = action.GetNonEmptyString("account_id")
	action.reference = action.GetString("reference")
}

func (action *CoreReferencesAction) checkAllowed() {
	action.IsAllowed(action.accountID)
}

func (action *CoreReferencesAction) loadRecords() {
	q := action.CoreQ().References().ForAccount(action.accountID)

	if action.reference != "" {
		q = q.ByReference(action.reference)
	}

	records, err := q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get References from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.records = records
}
