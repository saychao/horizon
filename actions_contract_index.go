package horizon

import (
	"github.com/saychao/horizon/db2"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"gitlab.com/tokend/regources"
)

// ContractIndexAction renders a page of contracts
// filters by startTime, endTime, disputing state,
// contractorID, customerID
type ContractIndexAction struct {
	Action
	PagingParams     db2.PageQuery
	StartTime        *int64
	EndTime          *int64
	Disputing        *bool
	Completed        *bool
	Source           string
	Counterparty     string
	EscrowID         string
	ContractNumber   string
	ContractsRecords []regources.Contract
	Page             hal.Page
}

// JSON is a method for actions.JSON
func (action *ContractIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *ContractIndexAction) checkAllowed() {
	action.IsAllowed(action.Source, action.EscrowID)
}

func (action *ContractIndexAction) loadParams() {
	action.ValidateCursorAsDefault()
	action.StartTime = action.GetOptionalInt64("start_time")
	action.EndTime = action.GetOptionalInt64("end_time")
	action.Disputing = action.GetOptionalBool("disputing")
	action.Completed = action.GetOptionalBool("completed")
	action.Counterparty = action.GetString("counterparty")
	action.Source = action.GetString("source")
	action.ContractNumber = action.GetString("contract_number")
	action.EscrowID = action.GetString("escrow")
	action.PagingParams = action.GetPageQuery()
	action.Page.Filters = map[string]string{
		"disputing":       action.GetString("disputing"),
		"completed":       action.GetString("completed"),
		"start_time":      action.GetString("start_time"),
		"end_time":        action.GetString("end_time"),
		"counterparty":    action.Counterparty,
		"source":          action.Source,
		"contract_number": action.ContractNumber,
		"escrow":          action.EscrowID,
	}
}

func (action *ContractIndexAction) loadRecords() {
	q := action.HistoryQ().Contracts()
	if action.StartTime != nil {
		q = q.ByStartTime(*action.StartTime)
	}
	if action.EndTime != nil {
		q = q.ByEndTime(*action.EndTime)
	}
	if action.Disputing != nil {
		q = q.ByDisputeState(*action.Disputing)
	}
	if action.Completed != nil {
		q = q.ByCompletedState(*action.Completed)
	}
	if action.Counterparty != "" {
		q = q.ByCounterpartyID(action.Counterparty)
	}
	if action.Source != "" {
		q = q.ByCounterpartyID(action.Source)
	}
	if action.EscrowID != "" {
		q = q.ByEscrowID(action.EscrowID)
	}
	if action.ContractNumber != "" {
		q = q.ByContractNumber(action.ContractNumber)
	}

	historyContracts, err := q.Page(action.PagingParams).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get contracts records")
		action.Err = &problem.ServerError
		return
	}

	for _, contract := range historyContracts {
		action.Page.Add(resource.PopulateContract(contract))
	}
}

func (action *ContractIndexAction) loadPage() {
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
