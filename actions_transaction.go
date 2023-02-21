package horizon

import (
	"github.com/SafeRE-IT/horizon/db2"
	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/render/sse"
	"github.com/SafeRE-IT/horizon/resource"
	"gitlab.com/tokend/regources"
)

// This file contains the actions:
//
// TransactionIndexAction: pages of transactions
// TransactionShowAction: single transaction by sequence, by hash or id

// TransactionIndexAction renders a page of ledger resources, identified by
// a normal page query.
type TransactionIndexAction struct {
	Action
	LedgerFilter  int32
	AccountFilter string
	BalanceFilter string
	PagingParams  db2.PageQuery
	Records       []history.Transaction
	MetaLedger    history.Ledger
	Page          hal.Page
}

// JSON is a method for actions.JSON
func (action *TransactionIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

// SSE is a method for actions.SSE
func (action *TransactionIndexAction) SSE(stream sse.Stream) {
	action.Setup(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
	)
	action.Do(
		func() {
			// we will reuse this variable in sse, so re-initializing is required
			action.Records = []history.Transaction{}
		},
		action.loadRecords,
		func() {
			records := action.Records[:]

			for _, record := range records {
				res := resource.PopulateTransaction(record)
				stream.Send(sse.Event{
					ID:   res.PagingToken(),
					Data: res,
				})
				action.PagingParams.Cursor = res.PagingToken()
			}
		},
	)
}

const (
	maxTxPagSize uint64 = 1000
)

func (action *TransactionIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.AccountFilter = action.GetString("account_id")
	action.LedgerFilter = action.GetInt32("ledger_id")

	action.PagingParams = action.GetPageQuery()

	limit := action.GetUInt64("limit")
	if limit > db2.MaxPageSize {
		action.PagingParams.Limit = limit
		if limit > maxTxPagSize {
			action.PagingParams.Limit = maxTxPagSize
		}
	}
}

func (action *TransactionIndexAction) loadRecords() {
	q := action.HistoryQ()
	txs := q.Transactions()

	switch {
	case action.AccountFilter != "":
		txs.ForAccount(action.AccountFilter)
	case action.LedgerFilter > 0:
		txs.ForLedger(action.LedgerFilter)
	}

	// memorize ledger sequence before select to prevent data race
	latestLedger := int32(action.App.historyLatestLedgerGauge.Value())

	err := txs.Page(action.PagingParams).Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to get transactions")
		action.Err = &problem.ServerError
		return
	}

	if uint64(len(action.Records)) == action.PagingParams.Limit {
		// we fetched full page, probably there is something ahead
		latestLedger = action.Records[len(action.Records)-1].LedgerSequence
	}

	// load ledger close time
	if err := action.HistoryQ().LedgerBySequence(&action.MetaLedger, latestLedger); err != nil {
		action.Log.WithError(err).Error("failed to get ledger")
		action.Err = &problem.ServerError
		return
	}
}

func (action *TransactionIndexAction) loadPage() {
	for _, record := range action.Records {
		res := resource.PopulateTransaction(record)
		action.Page.Add(res)
	}

	action.Page.Embedded.Meta = &hal.PageMeta{
		LatestLedger: &hal.LatestLedgerMeta{
			Sequence: action.MetaLedger.Sequence,
			ClosedAt: action.MetaLedger.ClosedAt,
		},
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}

func (action *TransactionIndexAction) checkAllowed() {
	action.IsAllowed(action.AccountFilter)
}

// TransactionShowAction renders a ledger found by its sequence number.
type TransactionShowAction struct {
	Action
	HashOrID string
	Record   history.Transaction
	Resource regources.Transaction
}

func (action *TransactionShowAction) loadParams() {
	action.HashOrID = action.GetString("id")
}

func (action *TransactionShowAction) loadRecord() {
	action.Err = action.HistoryQ().TransactionByHashOrID(&action.Record, action.HashOrID)
}

func (action *TransactionShowAction) loadResource() {
	action.Resource = resource.PopulateTransaction(action.Record)
}

// JSON is a method for actions.JSON
func (action *TransactionShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadResource,
		func() { hal.Render(action.W, action.Resource) },
	)
}

func (action *TransactionShowAction) checkAllowed() {
	action.IsAllowed("")
}
