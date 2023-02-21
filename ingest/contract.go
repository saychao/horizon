package ingest

import (
	"encoding/json"
	"github.com/saychao/horizon/db2"
	regources "gitlab.com/tokend/regources/generated"
	"time"

	"github.com/saychao/horizon/db2/history"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

func contractCreate(is *Session, entry *xdr.LedgerEntry) error {
	contract := convertContract(entry.Data.MustContract())

	err := is.Ingestion.Contracts(contract)
	if err != nil {
		return errors.Wrap(err, "failed to ingest contract", logan.F{
			"contract_id": contract.ID,
		})
	}

	return nil
}

func contractUpdate(is *Session, entry *xdr.LedgerEntry) error {
	contract := convertContract(entry.Data.MustContract())

	err := is.Ingestion.HistoryQ().Contracts().Update(contract)
	if err != nil {
		return errors.Wrap(err, "failed to update contract", logan.F{
			"contract_id": contract.ID,
		})
	}

	return nil
}

func contractDelete(is *Session, key *xdr.LedgerKey) error {
	manageContractOp := is.Cursor.Operation().Body.MustManageContractOp()
	if manageContractOp.ContractId != key.Contract.ContractId {
		return errors.New("Expected manage contract op contract id to be equal ledger key contract id")
	}
	contractID := int64(manageContractOp.ContractId)

	switch manageContractOp.Data.Action {
	case xdr.ManageContractActionConfirmCompleted:
		return is.updateInvoicesContractStates(contractID,
			history.ReviewableRequestStateApproved,
			int32(xdr.ContractStateCustomerConfirmed)|
				int32(xdr.ContractStateContractorConfirmed),
			[]history.ReviewableRequestState{history.ReviewableRequestStateApproved})
	case xdr.ManageContractActionResolveDispute:
		isRevert := manageContractOp.Data.IsRevert
		if *isRevert {
			return is.updateInvoicesContractStates(contractID,
				history.ReviewableRequestStatePermanentlyRejected,
				int32(xdr.ContractStateRevertingResolve),
				[]history.ReviewableRequestState{history.ReviewableRequestStatePending,
					history.ReviewableRequestStateWaitingForConfirmation})
		}

		err := is.Ingestion.HistoryQ().Contracts().AddState(contractID, int32(xdr.ContractStateNotRevertingResolve))
		if err != nil {
			return errors.Wrap(err, "failed to update contract state")
		}

		return is.updateNotRevertingContractInvoices(contractID)
	default:
		return errors.New("Unexpected case in delete contract")
	}
}

func (is *Session) updateInvoicesContractStates(
	contractID int64,
	invoiceState history.ReviewableRequestState,
	contractState int32,
	oldStates []history.ReviewableRequestState,
) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().
		UpdateInvoicesStates(invoiceState, oldStates, contractID)
	if err != nil {
		return errors.Wrap(err, "failed to update invoices states")
	}

	err = is.Ingestion.HistoryQ().Contracts().AddState(contractID, contractState)
	if err != nil {
		return errors.Wrap(err, "failed to update contract state")
	}

	return nil
}

func (is *Session) updateNotRevertingContractInvoices(contractID int64) error {
	err := is.Ingestion.HistoryQ().ReviewableRequests().UpdateInvoicesStates(
		history.ReviewableRequestStateApproved,
		[]history.ReviewableRequestState{history.ReviewableRequestStateWaitingForConfirmation},
		contractID)
	if err != nil {
		return errors.Wrap(err, "failed to update approved invoices", logan.F{
			"contract_id": contractID,
		})
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().UpdateInvoicesStates(
		history.ReviewableRequestStatePermanentlyRejected,
		[]history.ReviewableRequestState{history.ReviewableRequestStatePending},
		contractID)
	if err != nil {
		return errors.Wrap(err, "failed to update rejected invoices", logan.F{
			"contract_id": contractID,
		})
	}

	return nil
}

func convertContract(rawContract xdr.ContractEntry) history.Contract {
	var initialDetails regources.Details
	_ = json.Unmarshal([]byte(string(rawContract.InitialDetails)), &initialDetails)

	var customerDetails regources.Details
	_ = json.Unmarshal([]byte(string(rawContract.CustomerDetails)), &customerDetails)

	var invoices []int64
	for _, item := range rawContract.InvoiceRequestsIDs {
		invoices = append(invoices, int64(item))
	}

	return history.Contract{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(rawContract.ContractId),
		},
		Contractor:      rawContract.Contractor.Address(),
		Customer:        rawContract.Customer.Address(),
		Escrow:          rawContract.Escrow.Address(),
		StartTime:       time.Unix(int64(rawContract.StartTime), 0).UTC(),
		EndTime:         time.Unix(int64(rawContract.EndTime), 0).UTC(),
		InitialDetails:  initialDetails,
		CustomerDetails: customerDetails,
		Invoices:        invoices,
		State:           int32(rawContract.State),
	}
}

func (is *Session) processManageContract(op xdr.ManageContractOp, result xdr.ManageContractResult) error {
	if result.Code != xdr.ManageContractResultCodeSuccess {
		return nil
	}

	switch op.Data.Action {
	case xdr.ManageContractActionAddDetails:
		return is.addContractDetails(string(op.Data.MustDetails()), int64(op.ContractId))
	case xdr.ManageContractActionStartDispute:
		return is.addContractDispute(string(op.Data.MustDisputeReason()), int64(op.ContractId))
	}

	return nil
}

func (is *Session) addContractDetails(xdrDetails string, contractID int64) error {
	var details regources.Details
	err := json.Unmarshal([]byte(xdrDetails), &details)

	if err != nil {
		return errors.Wrap(err, "failed to unmarshal contract details", logan.F{
			"contract_id": contractID,
		})
	}

	source := is.Cursor.OperationSourceAccount()
	closeTime := is.Cursor.LedgerCloseTime()

	contractDetails := history.ContractDetails{
		ContractID: contractID,
		Details:    details,
		Author:     source.Address(),
		CreatedAt:  closeTime,
	}

	err = is.Ingestion.ContractDetails(contractDetails)
	if err != nil {
		return errors.Wrap(err, "failed to ingest contract details", logan.F{
			"contract_id": contractID,
		})
	}

	return nil
}

func (is *Session) addContractDispute(xdrReason string, contractID int64) error {
	var details regources.Details
	err := json.Unmarshal([]byte(xdrReason), &details)

	if err != nil {
		return errors.Wrap(err, "failed to unmarshal contract dispute reason", logan.F{
			"contract_id": contractID,
		})
	}

	source := is.Cursor.OperationSourceAccount()
	closeTime := is.Cursor.LedgerCloseTime()

	contractDispute := history.ContractDispute{
		ContractID: contractID,
		Reason:     details,
		Author:     source.Address(),
		CreatedAt:  closeTime,
	}

	err = is.Ingestion.ContractDispute(contractDispute)
	if err != nil {
		return errors.Wrap(err, "failed to ingest contract dispute", logan.F{
			"contract_id": contractID,
		})
	}

	return nil
}
