package resource

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

func PopulateKeyValue(keyValue *core.KeyValue) (*regources.KeyValue, error) {
	k := &regources.KeyValue{
		Key: keyValue.Key,
		Type: regources.Flag{
			Name:  keyValue.Value.Type.ShortString(),
			Value: int32(keyValue.Value.Type),
		},
	}

	switch keyValue.Value.Type {
	case xdr.KeyValueEntryTypeUint32:
		if keyValue.Value.Ui32Value != nil {
			uint32Value := uint32(*keyValue.Value.Ui32Value)
			k.Ui32Value = &uint32Value
		}
	case xdr.KeyValueEntryTypeUint64:
		if keyValue.Value.Ui64Value != nil {
			uint64Value := uint64(*keyValue.Value.Ui64Value)
			k.Ui64Value = &uint64Value
		}
	case xdr.KeyValueEntryTypeString:
		k.StringValue = keyValue.Value.StringValue
	default:
		return nil, errors.New("Unexpected key value type")
	}

	return k, nil
}
