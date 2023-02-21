package history

import (
	"github.com/SafeRE-IT/horizon/db2"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

var selectContracts = sq.Select(
	"hc.id",
	"hc.contractor",
	"hc.customer",
	"hc.escrow",
	"hc.start_time",
	"hc.end_time",
	"hc.initial_details",
	"hc.customer_details",
	"hc.invoices",
	"hc.state",
).From("history_contracts hc")

type ContractQI interface {
	// ByStartTime - filters contracts by start time
	ByStartTime(seconds int64) ContractQI
	// ByEndTime - filters contracts by end time
	ByEndTime(seconds int64) ContractQI
	// ByDisputeState - filters contracts by dispute state
	ByDisputeState(isDisputing bool) ContractQI
	// ByCompletedState - filters contracts by completed state
	ByCompletedState(isCompleted bool) ContractQI
	// ByCounterpartyID - filters contracts by contractor id
	ByCounterpartyID(counterpartyID string) ContractQI
	// ByContractNumber - filters contracts by contractor id
	ByContractNumber(contractNumber string) ContractQI
	// ByCustomerID - filters contracts by customer id
	ByEscrowID(escrowID string) ContractQI
	// Page - applies page params
	Page(page db2.PageQuery) ContractQI
	// Select - selects contract by specifics filters
	Select() ([]Contract, error)
	// ByID - get contract by contract id
	ByID(id int64) (*Contract, error)
	// Update - update contract
	Update(contract Contract) error
	AddState(contractID int64, stateToAdd int32) error
}

type ContractQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Contracts() ContractQI {
	return &ContractQ{
		parent: q,
		sql:    selectContracts,
	}
}

func (q *ContractQ) ByStartTime(seconds int64) ContractQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("start_time >= ?", time.Unix(seconds, 0).UTC())

	return q
}

func (q *ContractQ) ByEndTime(seconds int64) ContractQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("end_time >= ?", time.Unix(seconds, 0).UTC())

	return q
}

func (q *ContractQ) ByDisputeState(isDisputing bool) ContractQI {
	if q.Err != nil {
		return q
	}

	disputeState := int32(xdr.ContractStateDisputing)

	if isDisputing {
		q.sql = q.sql.Where("state & ? = ?", disputeState, disputeState)
	} else {
		q.sql = q.sql.Where("state & ? = 0", disputeState)
	}

	return q
}

func (q *ContractQ) ByCompletedState(isCompleted bool) ContractQI {
	if q.Err != nil {
		return q
	}

	bothCompleted := int32(xdr.ContractStateContractorConfirmed) | int32(xdr.ContractStateCustomerConfirmed)
	escrowResolve := int32(xdr.ContractStateRevertingResolve) | int32(xdr.ContractStateNotRevertingResolve)

	if isCompleted {
		q.sql = q.sql.Where("((state & ? = ?) or (state & ? != 0))", bothCompleted, bothCompleted, escrowResolve)
	} else {
		q.sql = q.sql.Where("((state & ? != ?) and (state & ? = 0))", bothCompleted, bothCompleted, escrowResolve)
	}

	return q
}

func (q *ContractQ) Page(page db2.PageQuery) ContractQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "id")
	return q
}

func (q *ContractQ) Select() ([]Contract, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []Contract
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

func (q *ContractQ) ByID(id int64) (*Contract, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.sql = q.sql.Where(sq.Eq{"id": id})

	var result Contract
	q.Err = q.parent.Get(&result, q.sql)
	return &result, q.Err
}

func (q *ContractQ) ByCounterpartyID(counterpartyID string) ContractQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("((contractor = ?) or (customer = ?))", counterpartyID, counterpartyID)

	return q
}

func (q *ContractQ) ByEscrowID(escrowID string) ContractQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"escrow": escrowID})

	return q
}

func (q *ContractQ) ByContractNumber(contractNumber string) ContractQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("initial_details->>'contract_number' = ?", contractNumber)

	return q
}

// Update - update contract using it's ID
func (q *ContractQ) Update(contract Contract) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("history_contracts").SetMap(map[string]interface{}{
		"contractor":       contract.Contractor,
		"customer":         contract.Customer,
		"escrow":           contract.Escrow,
		"start_time":       contract.StartTime,
		"end_time":         contract.EndTime,
		"initial_details":  contract.InitialDetails,
		"customer_details": contract.CustomerDetails,
		"invoices":         contract.Invoices,
		"state":            contract.State,
	}).Where("id = ?", contract.ID)

	err := q.parent.Exec(query)
	return err
}

func (q *ContractQ) AddState(contractID int64, stateToAdd int32) error {
	if q.Err != nil {
		return q.Err
	}

	query := "UPDATE history_contracts SET state = (state | ?) WHERE id = ?"

	err := q.parent.ExecRaw(query, stateToAdd, contractID)
	if err != nil {
		return errors.Wrap(err, "failed to execute contract raw")
	}

	return nil
}
