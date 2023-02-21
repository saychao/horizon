package horizon

import (
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	"github.com/saychao/horizon/resource/reviewablerequest"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// ContractShowAction represents singe contract by id
// with invoices requests
type ContractShowAction struct {
	Action
	ContractID     int64
	ContractRecord regources.Contract
}

// JSON is a method for actions.JSON
func (action *ContractShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecords,
		action.checkAllowed,
		func() {
			hal.Render(action.W, action.ContractRecord)
		},
	)
}

func (action *ContractShowAction) checkAllowed() {
	action.IsAllowed(action.ContractRecord.Contractor, action.ContractRecord.Customer, action.ContractRecord.Escrow)
}

func (action *ContractShowAction) loadParams() {
	action.ContractID = action.GetInt64("id")
}

func (action *ContractShowAction) loadRecords() {

	contract, err := action.HistoryQ().Contracts().ByID(action.ContractID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get contract record")
		action.Err = &problem.ServerError
		return
	}

	if contract == nil {
		action.Err = &problem.NotFound
		return
	}

	action.ContractRecord = resource.PopulateContract(*contract)

	invoices, err := action.HistoryQ().ReviewableRequests().ByIDs(contract.Invoices).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get invoices records")
		action.Err = &problem.ServerError
		return
	}

	for _, invoice := range invoices {
		res, err := reviewablerequest.PopulateReviewableRequest(&invoice)
		if err != nil {
			action.Log.WithError(err).Error("Failed to populate invoice request")
			action.Err = &problem.ServerError
			return
		}

		action.ContractRecord.Invoices = append(action.ContractRecord.Invoices, *res)
	}

	additionDetails, err := action.HistoryQ().ContractsDetails().ByContractID(contract.ID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get additional contract details records")
		action.Err = &problem.ServerError
		return
	}

	for _, detail := range additionDetails {
		res := resource.PopulateContractDetails(detail)
		action.ContractRecord.AdditionalDetails = append(action.ContractRecord.AdditionalDetails, res)
	}

	if (contract.State & int32(xdr.ContractStateDisputing)) == 0 {
		return
	}

	dispute, err := action.HistoryQ().ContractDispute().ByContractID(contract.ID)
	if (err != nil) || (dispute == nil) {
		action.Log.WithError(err).Error("Failed to get contract dispute records")
		action.Err = &problem.ServerError
		return
	}

	action.ContractRecord.DisputeReason = resource.PopulateContractDispute(*dispute)
}
