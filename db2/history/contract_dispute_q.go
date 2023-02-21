package history

import (
	sq "github.com/Masterminds/squirrel"
)

var selectContractDispute = sq.Select(
	"hcd.contract_id",
	"hcd.reason",
	"hcd.author",
	"hcd.created_at",
).From("history_contracts_disputes hcd")

type ContractDisputeQI interface {
	// ByContractID - get contract details by contract id
	ByContractID(contractID int64) (*ContractDispute, error)
}

type ContractDisputeQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) ContractDispute() ContractDisputeQI {
	return &ContractDisputeQ{
		parent: q,
		sql:    selectContractDispute,
	}
}

func (q *ContractDisputeQ) ByContractID(contractID int64) (*ContractDispute, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.sql = q.sql.Where(sq.Eq{"contract_id": contractID}).Limit(1)

	var result ContractDispute
	q.Err = q.parent.Get(&result, q.sql)
	return &result, q.Err
}
