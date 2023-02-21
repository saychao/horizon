package resources

import (
	"github.com/SafeRE-IT/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewKeyValue(record core2.KeyValue) regources.KeyValueEntry {
	return regources.KeyValueEntry{
		Key: regources.Key{
			Type: regources.KEY_VALUE_ENTRIES,
			ID:   record.Key,
		},
		Attributes: regources.KeyValueEntryAttributes{
			Value: regources.KeyValueEntryValue{
				Type: record.Value.Type,
				U32:  (*uint32)(record.Value.Ui32Value),
				U64:  (*uint64)(record.Value.Ui64Value),
				Str:  (*string)(record.Value.StringValue),
			},
		},
	}
}
