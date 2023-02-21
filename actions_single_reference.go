package horizon

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
)

type SingleReferenceAction struct {
	Action

	accountID string
	reference string

	record core.Reference
}

func (action *SingleReferenceAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		func() {
			response := map[string]interface{}{
				"data": action.record,
			}
			hal.Render(action.W, response)
		},
	)
}

func (action *SingleReferenceAction) loadParams() {
	action.accountID = action.GetNonEmptyString("account_id")
	action.reference = action.GetNonEmptyString("reference")
}

func (action *SingleReferenceAction) loadRecords() {
	q := action.CoreQ().References().
		ForAccount(action.accountID).
		ByReference(action.reference)

	records, err := q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get References from core DB")
		action.Err = &problem.ServerError
		return
	}

	if len(records) == 0 {
		action.Err = &problem.NotFound
		return
	}

	action.record = records[0]
}
