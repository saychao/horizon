package ingest

import (
	"encoding/json"

	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/utf8"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

func (is *Session) processReviewRequest(op xdr.ReviewRequestOp, res xdr.ExtendedResult,
	changes xdr.LedgerEntryChanges) (err error) {

	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err = is.approveReviewableRequest(op, res, changes)
	case xdr.ReviewRequestOpActionPermanentReject:
		err = is.permanentReject(op)
	case xdr.ReviewRequestOpActionReject:
		return
	default:
		err = errors.From(errors.New("Unexpected review request action"), map[string]interface{}{
			"action_type": op.Action,
		})
	}

	if err != nil {
		return errors.Wrap(err, "failed to process review request", map[string]interface{}{
			"request_id": uint64(op.RequestId),
		})
	}
	return nil
}

func hasDeletedReviewableRequest(changes xdr.LedgerEntryChanges) bool {
	for i := range changes {
		if changes[i].Removed == nil {
			continue
		}

		if changes[i].Removed.ReviewableRequest != nil {
			return true
		}
	}

	return false
}

func (is *Session) approveReviewableRequest(op xdr.ReviewRequestOp, res xdr.ExtendedResult,
	changes xdr.LedgerEntryChanges) error {

	if op.RequestDetails.RequestType == xdr.ReviewableRequestTypeChangeRole && !hasDeletedReviewableRequest(changes) {
		return nil
	}

	if !res.Fulfilled {
		return nil
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().Approve(uint64(op.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to approve reviewable request")
	}

	switch op.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeCreateWithdraw:
		err = is.setWithdrawalDetails(uint64(op.RequestId), op.RequestDetails.Withdrawal)
	case xdr.ReviewableRequestTypeCreateInvoice:
		err = is.setWaitingForConfirmationState(uint64(op.RequestId))
	}

	if err != nil {
		return errors.Wrap(err, "failed to set request details")
	}

	return nil
}

func (is *Session) setWithdrawalDetails(requestID uint64, details *xdr.WithdrawalDetails) error {
	fields := logan.Field("request_id", requestID)
	request, err := is.Ingestion.HistoryQ().ReviewableRequests().ByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to load reviewable request by id", fields)
	}

	if request == nil {
		return errors.From(errors.New("reviewable request not found"), fields)
	}

	if request.RequestType != xdr.ReviewableRequestTypeCreateWithdraw {
		return errors.From(errors.New("expected withdrawal request"), fields.Add("request_type", request.RequestType))
	}

	var reviewerDetails map[string]interface{}

	externalDetails := utf8.Scrub(details.ExternalDetails)
	err = json.Unmarshal([]byte(externalDetails), &reviewerDetails)
	if err != nil {
		// we ignore here error on purpose, as it's too late to valid the error
		err = errors.Wrap(err, "failed to marshal reviewer details", fields)
		is.log.WithError(err).WithFields(logan.F{
			"scrubbed_details": externalDetails,
			"original_details": details.ExternalDetails,
		}).Warn("Reviewer sent invalid json in withdrawal details")
	}

	var withdrawalDetails *history.WithdrawalRequest
	if request.Details.Withdraw != nil {
		withdrawalDetails = request.Details.Withdraw
	} else if request.Details.TwoStepWithdraw != nil {
		withdrawalDetails = request.Details.TwoStepWithdraw
	} else {
		return errors.New("Unexpected state: expected withdrawal details to be available")
	}

	withdrawalDetails.ReviewerDetails = reviewerDetails
	err = is.Ingestion.HistoryQ().ReviewableRequests().Update(*request)
	if err != nil {
		return errors.Wrap(err, "failed to update withdrawal request", fields)
	}

	return nil
}

func (is *Session) setWaitingForConfirmationState(requestID uint64) error {
	request, err := is.Ingestion.HistoryQ().ReviewableRequests().ByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to get request", logan.F{
			"request_id": requestID,
		})
	}

	if (request == nil) || (request.Details.Invoice == nil) || (request.Details.Invoice.ContractID == nil) {
		return nil
	}

	err = is.Ingestion.HistoryQ().ReviewableRequests().UpdateStates([]int64{int64(requestID)},
		history.ReviewableRequestStateWaitingForConfirmation)
	if err != nil {
		return errors.Wrap(err, "failed to update request state")
	}

	return nil
}
