//Package codes is a helper package to help convert to transaction and operation result codes
//to strings used in horizon.
package codes

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// ErrUnknownCode is returned when an unexepcted value is provided to `String`
var ErrUnknownCode = errors.New("Unknown result code")

type shortStr interface {
	ShortString() string
}

//opCodeToString returns the appropriate string representation of the provided result code
func opCodeToString(codeProvider shortStr) string {
	return "op_" + codeProvider.ShortString()
}

// ForOperationResult returns the strong represtation used by horizon for the
// error code `opr`
func ForOperationResult(opr xdr.OperationResult) (string, string, error) {
	if opr.Code != xdr.OperationResultCodeOpInner {
		return opr.Code.ShortString(), getMessage(opr.Code.ShortString()), nil
	}

	ir := opr.MustTr()
	ic, ok := codeProviders[ir.Type]
	if !ok {
		return "", "", errors.Wrap(ErrUnknownCode, "failed to get code provider", map[string]interface{}{
			"type": ir.Type.String(),
		})
	}

	code := ic(ir)
	opCode := opCodeToString(code)
	return opCode, getMessage(opCode), nil
}

func ForTxResult(txResult xdr.TransactionResult) (txResultCode string, opResultCodes []string, messages []string, err error) {
	txResultCode = txResult.Result.Code.ShortString()

	if txResult.Result.Results == nil {
		opResultCodes = []string{getMessage(txResultCode)}
		return
	}

	opResults := txResult.Result.MustResults()
	opResultCodes = make([]string, len(opResults))
	messages = make([]string, len(opResults))
	for i := range opResults {
		opResultCodes[i], messages[i], err = ForOperationResult(opResults[i])
		if err != nil {
			err = errors.Wrap(err, "Failed to convert to string op result")
			return
		}
	}

	return
}
