package horizon

import (
	"github.com/saychao/horizon/db2/core"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
)

type KeyValueShowAction struct {
	Action
	key            string
	keyValueRecord *core.KeyValue
}

func (action *KeyValueShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			res, err := resource.PopulateKeyValue(action.keyValueRecord)
			if err != nil {
				action.Log.WithError(err).Error("Failed to populate key_value")
				action.Err = &problem.ServerError
				return
			}
			hal.Render(action.W, *res)
		},
	)
}

func (action *KeyValueShowAction) loadParams() {
	action.key = action.GetString("key")
}

func (action *KeyValueShowAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.keyValueRecord, err = q.ByKey(action.key)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get key_value from core DB")
		action.Err = &problem.ServerError
		return
	}

	if action.keyValueRecord == nil {
		action.Err = &problem.NotFound
		return
	}
}

type KeyValueShowAllAction struct {
	Action
	coreRecords     []core.KeyValue
	recordsToRender []regources.KeyValue
}

func (action *KeyValueShowAllAction) JSON() {
	action.Do(
		action.loadRecord,
		func() {
			hal.Render(action.W, action.recordsToRender)
		},
	)
}

func (action *KeyValueShowAllAction) loadRecord() {
	var err error
	q := action.CoreQ().KeyValue()
	action.coreRecords, err = q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get all key_values from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.recordsToRender, err = action.getPopulatedKeyValues()
	if err != nil {
		action.Log.WithError(err).Error("Failed to populate all key_values")
		action.Err = &problem.ServerError
		return
	}
}

func (action *KeyValueShowAllAction) getPopulatedKeyValues() ([]regources.KeyValue, error) {
	var res []regources.KeyValue
	for _, keyValue := range action.coreRecords {
		keyValueEntry, err := resource.PopulateKeyValue(&keyValue)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to populate key_value", logan.F{
				"keyValue": keyValue,
			})
		}
		res = append(res, *keyValueEntry)
	}
	return res, nil
}
