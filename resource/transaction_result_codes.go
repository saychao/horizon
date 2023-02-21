package resource

import (
	"github.com/SafeRE-IT/horizon/codes"
	"gitlab.com/tokend/go/xdr"
)

// TransactionResultCodes represent a summary of result codes returned from
// a single xdr TransactionResult
type TransactionResultCodes struct {
	TransactionCode string   `json:"transaction"`
	OperationCodes  []string `json:"operations,omitempty"`
	Messages        []string `json:"messages"`
}

// NewTransactionResultCodes creates tx result codes from xdr result
func NewTransactionResultCodes(txResult xdr.TransactionResult) (*TransactionResultCodes, error) {
	txResultCode, opResultCodes, messages, err := codes.ForTxResult(txResult)
	if err != nil {
		return nil, err
	}

	return &TransactionResultCodes{
		TransactionCode: txResultCode,
		OperationCodes:  opResultCodes,
		Messages:        messages,
	}, nil
}
