package horizon

import (
	"github.com/SafeRE-IT/horizon/db2"
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
	"gitlab.com/tokend/go/xdr"
)

// This file contains the actions:
//
// OperationIndexAction: pages of operations
// OperationShowAction: single operation by id

// OperationIndexAction renders a page of operations resources, identified by
// a normal page query and optionally filtered by an account, ledger, or
// transaction.
type HistoryOperationIndexAction struct {
	Action
	Types        []xdr.OperationType
	PagingParams db2.PageQuery
	Records      []history.Operation
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *HistoryOperationIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *HistoryOperationIndexAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
}

func (action *HistoryOperationIndexAction) loadRecords() {
	err := action.HistoryQ().Operations().Page(action.PagingParams).Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to get operations")
		action.Err = &problem.ServerError
		return
	}
}

func (action *HistoryOperationIndexAction) loadPage() {
	for _, record := range action.Records {
		var res hal.Pageable

		res, action.Err = resource.NewPublicOperation(action.Ctx, record, nil)
		if action.Err != nil {
			return
		}
		action.Page.Add(res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

// OperationShowAction renders a ledger found by its sequence number.
type HistoryOperationShowAction struct {
	Action
	ID           int64
	Record       history.Operation
	Resource     interface{}
	Participants map[int64]*history.OperationParticipants
}

func (action *HistoryOperationShowAction) loadParams() {
	action.ID = action.GetInt64("id")
}

func (action *HistoryOperationShowAction) loadRecord() {
	action.Err = action.HistoryQ().OperationByID(&action.Record, action.ID)
	if action.Err != nil {
		return
	}

	action.Participants = map[int64]*history.OperationParticipants{
		action.ID: {},
	}
	switch action.Record.Type {
	case xdr.OperationTypeManageOffer:
		// workaround for load participants
		action.IsAdmin = true
		action.LoadParticipants("", action.Participants)
		// reverting workaround, just in case
		action.IsAdmin = false
	}
}

func (action *HistoryOperationShowAction) loadResource() {
	action.Resource, action.Err = resource.NewPublicOperation(
		action.Ctx, action.Record, action.Participants[action.Record.ID].Participants)
}

// JSON is a method for actions.JSON
func (action *HistoryOperationShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}
