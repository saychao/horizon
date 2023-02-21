package horizon

import (
	"fmt"
	"net/http"

	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/render/problem"
	"github.com/saychao/horizon/resource"
	txsubHelper "github.com/saychao/horizon/txsub"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/txsub"
	"gitlab.com/tokend/go/xdr"
)

// TransactionCreateAction submits a transaction to the stellar-core network
// on behalf of the requesting client.
type TransactionCreateAction struct {
	Action
	tx       string
	result   txsub.Result
	resource resource.TransactionSuccess
}

// JSON format action handler
func (action *TransactionCreateAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadTX,
		action.loadResult,
		action.loadResource,
		func() {
			hal.Render(action.W, action.resource)
		})
}

func (action *TransactionCreateAction) loadTX() {
	action.tx = action.GetNonEmptyString("tx")
}

func (action *TransactionCreateAction) loadResult() {
	envelopeInfo, err := txsubHelper.ExtractEnvelopeInfo(action.tx, action.App.CoreInfo.NetworkPassphrase)
	if err != nil {
		action.Err = &problem.P{
			Type:   "transaction_malformed",
			Title:  "Transaction Malformed",
			Status: http.StatusBadRequest,
			Detail: "Horizon could not decode the transaction envelope in this " +
				"request. A transaction should be an XDR TransactionEnvelope struct " +
				"encoded using base64.  The envelope read from this request is " +
				"echoed in the `extras.envelope_xdr` field of this response for your " +
				"convenience.",
			Extras: map[string]interface{}{
				"envelope_xdr": action.tx,
			},
		}
		return
	}

	action.result = action.App.submitter.Submit(action.Ctx, envelopeInfo)
	if action.result.HasInternalError() {
		action.Log.WithError(action.result.Err).Error("Failed to submit tx")
		action.Err = &problem.ServerError
		return
	}

	if action.result.Err == nil {
		action.resource.Populate(action.Ctx, action.result)
		return
	}
}

func (action *TransactionCreateAction) loadResource() {
	p, err := txResultToProblem(&action.result)
	if err != nil {
		action.Log.WithError(err).Error("failed to craft problem")
		action.Err = &problem.ServerError
		return
	}

	if p != nil {
		action.Err = p
		return
	}
}

func txResultToProblem(result *txsub.Result) (*problem.P, error) {
	if result.Err == nil {
		return nil, nil
	}

	txSubError, ok := result.Err.(txsub.Error)
	if !ok {
		return nil, errors.New("Unexpected error type")
	}

	switch txSubError.Type() {
	case txsub.Timeout:
		return &problem.Timeout, nil
	case txsub.RejectedTx:
		var parsedResult xdr.TransactionResult
		err := xdr.SafeUnmarshalBase64(txSubError.ResultXDR(), &parsedResult)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get parse tx result")
		}

		var parsedEnvelope xdr.TransactionEnvelope
		err = xdr.SafeUnmarshalBase64(result.EnvelopeXDR, &parsedEnvelope)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal tx envelope")
		}

		resultCodes, err := resource.NewTransactionResultCodes(parsedResult)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create transaction result codes")
		}

		return &problem.P{
			Type:   "transaction_failed",
			Title:  "Transaction Failed",
			Status: http.StatusBadRequest,
			Detail: "The transaction failed when submitted to the stellar network. " +
				"The `extras.result_codes` field on this response contains further " +
				"details.  Descriptions of each code can be found at: " +
				"https://www.stellar.org/developers/learn/concepts/list-of-operations.html",
			Extras: map[string]interface{}{
				"envelope_xdr":    result.EnvelopeXDR,
				"result_xdr":      txSubError.ResultXDR(),
				"result_codes":    resultCodes,
				"parsed_result":   &parsedResult,
				"parsed_envelope": &parsedEnvelope,
			},
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unexpected error type: %d", txSubError.Type()))
	}
}
