package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokend/go/strkey"
)

func newBalanceID(r *base, name string) (string, error) {
	rawID := r.getString(name)
	_, err := strkey.Decode(strkey.VersionByteBalanceID, rawID)
	if err != nil {
		return "", validation.Errors{
			name: err,
		}
	}
	return rawID, nil
}
