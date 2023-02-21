package resources

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SafeRE-IT/horizon/resource"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"

	"github.com/google/jsonapi"

	"github.com/SafeRE-IT/horizon/txsub/v2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewTxKey - creates new Tx Key for specified ID
func NewTxKey(txID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(txID, 10),
		Type: regources.TRANSACTIONS,
	}
}

func NewTxKeyFromHash(txHash string) regources.Key {
	return regources.Key{
		ID:   txHash,
		Type: regources.TRANSACTIONS,
	}
}

func NewTxFailure(env txsub.EnvelopeInfo, txSubErr txsub.Error) error {
	if txSubErr.Type() != txsub.RejectedTx { // timeout
		return &jsonapi.ErrorObject{
			Title:  http.StatusText(txSubErr.Code()),
			Detail: txSubErr.Details(),
			Status: fmt.Sprintf("%d", txSubErr.Code()),
		}
	}

	var parsedResult xdr.TransactionResult
	err := xdr.SafeUnmarshalBase64(txSubErr.ResultXDR(), &parsedResult)
	if err != nil {
		return errors.Wrap(err, "Failed to get parse tx result")
	}

	var parsedEnvelope xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(env.RawBlob, &parsedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal tx envelope")
	}
	resultCodes, err := resource.NewTransactionResultCodes(parsedResult)
	if err != nil {
		return errors.Wrap(err, "Failed to create transaction result codes")
	}
	meta := map[string]interface{}{
		"envelope":        env.RawBlob,
		"result_xdr":      txSubErr.ResultXDR(),
		"parsed_envelope": parsedEnvelope,
		"parsed_result":   parsedResult,
		"result_codes":    resultCodes,
	}

	return &jsonapi.ErrorObject{
		Title:  http.StatusText(txSubErr.Code()),
		Detail: txSubErr.Details(),
		Status: fmt.Sprintf("%d", txSubErr.Code()),
		Meta:   &meta,
	}
}
