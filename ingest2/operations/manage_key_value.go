package operations

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type manageKeyValueOpHandler struct {
	effectsProvider
}

// Details returns details about manage key value operation
func (h *manageKeyValueOpHandler) Details(op rawOperation, _ xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	manageKVOp := op.Body.MustManageKeyValueOp()

	var value *regources.KeyValueEntryValue
	if manageKVOp.Action.Action == xdr.ManageKvActionPut {
		valueForPtr := manageKVOp.Action.MustValue()
		value = &regources.KeyValueEntryValue{
			Type: valueForPtr.Type,
		}

		switch valueForPtr.Type {
		case xdr.KeyValueEntryTypeUint32:
			value.U32 = new(uint32)
			*value.U32 = uint32(*valueForPtr.Ui32Value)
		case xdr.KeyValueEntryTypeString:
			value.Str = new(string)
			*value.Str = string(*valueForPtr.StringValue)
		case xdr.KeyValueEntryTypeUint64:
			value.U64 = new(uint64)
			*value.U64 = uint64(*valueForPtr.Ui64Value)
		default:
			return history2.OperationDetails{}, errors.From(errors.New("unexpected key value value type"), logan.F{
				"type": valueForPtr.Type.ShortString(),
			})
		}
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeManageKeyValue,
		ManageKeyValue: &history2.ManageKeyValueDetails{
			Key:    string(manageKVOp.Key),
			Action: manageKVOp.Action.Action,
			Value:  value,
		},
	}, nil
}
