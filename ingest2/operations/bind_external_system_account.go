package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"gitlab.com/tokend/go/xdr"
)

type bindExternalSystemAccountOpHandler struct {
	effectsProvider
}

// Details returns details about bind external system account operation
func (h *bindExternalSystemAccountOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	bindExternalSystemAccountOp := op.Body.MustBindExternalSystemAccountIdOp()

	return history2.OperationDetails{
		Type: xdr.OperationTypeBindExternalSystemAccountId,
		BindExternalSystemAccount: &history2.BindExternalSystemAccountDetails{
			ExternalSystemType: int32(bindExternalSystemAccountOp.ExternalSystemType),
		},
	}, nil
}
