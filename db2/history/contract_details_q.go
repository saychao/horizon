package history

import (
	sq "github.com/Masterminds/squirrel"
)

var selectContractsDetails = sq.Select(
	"hcd.contract_id",
	"hcd.details",
	"hcd.author",
	"hcd.created_at",
).From("history_contracts_details hcd")

type ContractsDetailsQI interface {
	// ByContractID - get contract details by contract id
	ByContractID(contractID int64) ([]ContractDetails, error)
}

type ContractsDetailsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) ContractsDetails() ContractsDetailsQI {
	return &ContractsDetailsQ{
		parent: q,
		sql:    selectContractsDetails,
	}
}

func (q *ContractsDetailsQ) ByContractID(contractID int64) ([]ContractDetails, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.sql = q.sql.Where(sq.Eq{"contract_id": contractID})

	var result []ContractDetails
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}
