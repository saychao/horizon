package core2

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type KeyValueEntry xdr.KeyValueEntryValue

func (k *KeyValueEntry) Scan(src interface{}) error {
	var data string
	switch rawData := src.(type) {
	case []byte:
		data = string(rawData)
	case string:
		data = rawData
	default:
		return errors.New("Unexpected type for KeyValueEntry")
	}

	var rawData xdr.KeyValueEntryValue
	err := xdr.SafeUnmarshalBase64(data, &rawData)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal key_value")
	}

	*k = KeyValueEntry(rawData)

	return nil
}

type KeyValue struct {
	Key   string        `db:"key"`
	Value KeyValueEntry `db:"value"`
}
