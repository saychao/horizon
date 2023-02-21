package txsub

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ErrorType defines specific type of error
//go:generate jsonenums -type=ErrorType
type ErrorType int

const (
	// Timeout error during submission
	Timeout ErrorType = iota
	// RejectedTx occurs when Tx was reject by core
	RejectedTx
	// Canceled error cause of context done
	Canceled
)

var (
	txSubCodes = map[ErrorType]int{
		Timeout:    http.StatusRequestTimeout,
		RejectedTx: http.StatusBadRequest,
		Canceled:   http.StatusGone,
	}
	txSubDetails = map[ErrorType]string{
		Timeout:    "Tx submit is taking too long",
		RejectedTx: "Tx is invalid and got rejected",
		Canceled:   "User canceled request",
	}
)

// An Error represents a transaction submission error.
type Error interface {
	error
	// Type -- returns type of the error
	Type() ErrorType
	// ResultXDR returns base64 xdr encoded xdr.TransactionResult
	// Returns empty string for all non rejected txs
	ResultXDR() string

	Code() int

	Details() string
}

var timeoutError = &txSubError{
	error:     errors.New("timeout"),
	errorType: Timeout,
}

var cancelError = &txSubError{
	error:     errors.New("canceled"),
	errorType: Canceled,
}

type txSubError struct {
	error
	errorType ErrorType
	resultXDR string
}

func (m *txSubError) Type() ErrorType {
	return m.errorType
}

func (m *txSubError) ResultXDR() string {
	return m.resultXDR
}

func (m *txSubError) Code() int {
	return txSubCodes[m.errorType]
}

func (m *txSubError) Details() string {
	return txSubDetails[m.errorType]
}

// NewRejectedTxError creates error with type RejectedTx
func NewRejectedTxError(resultXDR string) Error {

	return &txSubError{
		error: errors.From(errors.New("core rejected tx"), logan.F{
			"result_xdr": resultXDR,
		}),
		errorType: RejectedTx,
		resultXDR: resultXDR,
	}
}

// IsInternalError - returns true if error is internal
func IsInternalError(err error) bool {
	if err == nil {
		return false
	}

	_, isTxError := err.(Error)
	return isTxError
}
