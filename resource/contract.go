package resource

import (
	"encoding/json"
	"strconv"

	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/resource/base"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateContract(row history.Contract) regources.Contract {
	var initialDetails map[string]interface{}
	_ = json.Unmarshal(row.InitialDetails, &initialDetails)

	var customerDetails map[string]interface{}
	_ = json.Unmarshal(row.CustomerDetails, &customerDetails)

	return regources.Contract{
		ID:              strconv.FormatInt(row.ID, 10),
		PT:              row.PagingToken(),
		Contractor:      row.Contractor,
		Customer:        row.Customer,
		Escrow:          row.Escrow,
		StartTime:       row.StartTime,
		EndTime:         row.EndTime,
		InitialDetails:  initialDetails,
		CustomerDetails: customerDetails,
		State:           base.FlagFromXdrContractState(row.State, xdr.ContractStateAll),
	}
}

func PopulateContractDetails(row history.ContractDetails) regources.DetailsWithPayload {
	var details map[string]interface{}
	_ = json.Unmarshal(row.Details, &details)

	return regources.DetailsWithPayload{
		Details:   details,
		Author:    row.Author,
		CreatedAt: row.CreatedAt,
	}
}

func PopulateContractDispute(row history.ContractDispute) *regources.DetailsWithPayload {
	var reason map[string]interface{}
	_ = json.Unmarshal(row.Reason, &reason)

	return &regources.DetailsWithPayload{
		Details:   reason,
		Author:    row.Author,
		CreatedAt: row.CreatedAt,
	}
}
