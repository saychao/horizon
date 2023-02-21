package horizon

import (
	"fmt"
	"github.com/saychao/horizon/db2"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/ledger"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/render/sse"
	"github.com/saychao/horizon/resource"
	"github.com/saychao/horizon/toid"
	"gitlab.com/tokend/go/xdr"
)

// This file contains the actions:
//
// OperationIndexAction: pages of operations
// OperationShowAction: single operation by id

// OperationIndexAction renders a page of operations resources, identified by
// a normal page query and optionally filtered by an account, ledger, or
// transaction.
type OperationIndexAction struct {
	Action
	Types []xdr.OperationType

	LedgerFilter        int32
	AccountFilter       string
	AccountTypeFilter   int32
	BalanceFilter       string
	AssetFilter         string
	ExchangeFilter      string
	TransactionFilter   string
	CompletedOnlyFilter bool
	SkipCanceled        bool
	PendingOnlyFilter   bool
	// ReferenceFilter substring
	ReferenceFilter string
	SinceFilter     *time.Time
	ToFilter        *time.Time
	PagingParams    db2.PageQuery
	Records         []history.Operation
	Participants    map[int64]*history.OperationParticipants
	Page            hal.Page
}

// JSON is a method for actions.JSON
func (action *OperationIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadParticipants,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		})
}

// SSE is a method for actions.SSE
func (action *OperationIndexAction) SSE(stream sse.Stream) {
	action.Setup(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
	)
	action.Do(
		func() {
			// we will reuse this variable in sse, so re-initializing is required
			action.Records = []history.Operation{}
		},
		action.loadRecords,
		action.loadParticipants,
		func() {
			records := action.Records[:]

			for _, record := range records {
				opParticipants := action.Participants[record.ID]
				res, err := resource.NewOperation(action.Ctx, record, opParticipants.Participants)

				if err != nil {
					stream.Err(action.Err)
					return
				}

				stream.Send(sse.Event{
					ID:   res.PagingToken(),
					Data: res,
				})
				action.PagingParams.Cursor = res.PagingToken()
			}
		})

}

func (action *OperationIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.AccountFilter = action.GetString("account_id")
	action.AccountTypeFilter = action.GetInt32("account_type")
	action.BalanceFilter = action.GetString("balance_id")
	action.AssetFilter = action.GetString("asset")
	action.TransactionFilter = action.GetString("tx_id")
	action.ReferenceFilter = action.GetString("reference")
	action.SinceFilter = action.TryGetTime("since")
	action.ToFilter = action.TryGetTime("to")
	action.CompletedOnlyFilter = action.GetBoolOrDefault("completed_only", true)
	action.SkipCanceled = action.GetBoolOrDefault("skip_canceled", true)
	action.PendingOnlyFilter = action.GetBool("pending_only")

	if action.CompletedOnlyFilter && action.PendingOnlyFilter {
		action.SetInvalidField("pending_only", errors.New("completed_only and pending_only filters cannot both be set"))
		return
	}

	var err error
	opTypeStr := action.GetString("operation_type")
	opType := int64(0)
	if opTypeStr != "" {
		opType, err = strconv.ParseInt(opTypeStr, 10, 64)
		if err != nil {
			action.SetInvalidField("operation_type", err)
		}
	}
	if opTypeStr != "" && len(action.Types) > 0 {
		// operations were set already, so action is limited to some types
		// checking if specified is one of them
		opType = func() int64 {
			for _, t := range action.Types {
				if int64(t) == opType {
					return opType
				}
			}
			return 0
		}()

		if opType == 0 {
			action.SetInvalidField(
				"operation_type", errors.New("invalid in some way"))
		}
	}
	if opTypeStr != "" {
		action.Types = []xdr.OperationType{xdr.OperationType(opType)}
	}

	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"account_id":     action.AccountFilter,
		"account_type":   fmt.Sprintf("%d", action.AccountTypeFilter),
		"balance_id":     action.BalanceFilter,
		"asset":          action.AssetFilter,
		"tx_id":          action.TransactionFilter,
		"reference":      action.ReferenceFilter,
		"since":          action.GetString("since"),
		"to":             action.GetString("to"),
		"completed_only": action.GetString("completed_only"),
		"skip_canceled":  action.GetString("skip_canceled"),
		"pending_only":   action.GetString("pending_only"),
	}
	if action.SinceFilter != nil {
		action.Page.Filters["since"] = action.SinceFilter.Format(time.RFC3339)
	}
	if action.ToFilter != nil {
		action.Page.Filters["to"] = action.ToFilter.Format(time.RFC3339)
	}
	if opType != 0 {
		action.Page.Filters["operation_type"] = fmt.Sprintf("%d", opType)
	}
}

func (action *OperationIndexAction) loadRecords() {
	ops := action.HistoryQ().Operations().WithoutCancelingManagerOffer().WithoutExternallyFullyMatched()

	if len(action.Types) > 0 {
		ops.ForTypes(action.Types)
	}

	if action.SkipCanceled {
		ops = ops.WithoutCanceled()
	}

	if action.AccountFilter != "" || action.AccountTypeFilter != 0 {
		ops.JoinOnAccount()
	}

	if action.AccountFilter != "" {
		ops.ForAccount(action.AccountFilter)
	}

	if action.AccountTypeFilter != 0 {
		ops.ForAccountType(action.AccountTypeFilter)
	}

	if action.BalanceFilter != "" {
		ops.ForBalance(action.BalanceFilter)
	}

	if action.AssetFilter != "" || action.ExchangeFilter != "" {
		ops.JoinOnBalance()
	}

	if action.AssetFilter != "" {
		ops.ForAsset(action.AssetFilter)
	}

	if action.TransactionFilter != "" {
		ops.ForTx(action.TransactionFilter)
	}

	if action.ReferenceFilter != "" {
		ops.ForReference(action.ReferenceFilter)
	}

	if action.SinceFilter != nil {
		ops.Since(*action.SinceFilter)
	}

	if action.ToFilter != nil {
		ops.To(*action.ToFilter)
	}

	if action.CompletedOnlyFilter {
		ops.CompletedOnly()
	}

	if action.PendingOnlyFilter {
		ops.PendingOnly()
	}

	err := ops.Page(action.PagingParams).Select(&action.Records)

	if err != nil {
		action.Log.WithError(err).Error("Failed to get operations")
		action.Err = &problem.ServerError
		return
	}
}

func (action *OperationIndexAction) loadParticipants() {
	// initializing our operation -> participants map
	action.Participants = map[int64]*history.OperationParticipants{}
	for _, operation := range action.Records {
		action.Participants[operation.ID] = &history.OperationParticipants{
			operation.Type,
			[]*history.Participant{},
		}
	}

	action.LoadParticipants(action.AccountFilter, action.Participants)
}

func (action *OperationIndexAction) loadPage() {
	for _, record := range action.Records {
		var res hal.Pageable
		opParticipants := action.Participants[record.ID]

		res, err := resource.NewOperation(action.Ctx, record, opParticipants.Participants)
		if err != nil {
			action.Log.WithError(err).Error("failed to populate operation")
			action.Err = &problem.ServerError
			return
		}

		// add operation 2 time if it's within one account
		if len(opParticipants.Participants) == 2 && record.Type != xdr.OperationTypeManageOffer && record.Type != xdr.OperationTypeCreateIssuanceRequest &&
			opParticipants.Participants[0].AccountID == opParticipants.Participants[1].AccountID {

			action.Page.Add(res)
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

func (action *OperationIndexAction) checkAllowed() {
	action.IsAllowed(action.AccountFilter)
}

// OperationShowAction renders a ledger found by its sequence number.
type OperationShowAction struct {
	Action
	ID           int64
	Record       history.Operation
	Participants map[int64]*history.OperationParticipants
	Resource     interface{}
}

func (action *OperationShowAction) loadParams() {
	action.ID = action.GetInt64("id")
}

func (action *OperationShowAction) loadRecord() {
	err := action.HistoryQ().OperationByID(&action.Record, action.ID)
	if err != nil {
		action.Log.WithError(err).Error("failed to load record")
		action.Err = &problem.ServerError
		return
	}

	action.Participants = map[int64]*history.OperationParticipants{
		action.ID: {},
	}
	action.LoadParticipants("", action.Participants)
}

func (action *OperationShowAction) loadResource() {
	action.Resource, action.Err = resource.NewOperation(action.Ctx, action.Record, action.Participants[action.Record.ID].Participants)
}

// JSON is a method for actions.JSON
func (action *OperationShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.verifyWithinHistory,
		action.loadRecord,
		action.loadResource,
	)
	action.Do(func() {
		hal.Render(action.W, action.Resource)
	})
}

func (action *OperationShowAction) verifyWithinHistory() {
	parsed := toid.Parse(action.ID)
	if parsed.LedgerSequence < ledger.CurrentState().History.OldestOnStart {
		action.Err = &problem.BeforeHistory
	}
}

func (action *OperationShowAction) checkAllowed() {
	action.IsAllowed("")
}
