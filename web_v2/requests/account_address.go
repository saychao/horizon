package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokend/go/strkey"
)

func newAccountAddress(r *base, name string) (string, error) {
	rawAddress := r.getString(name)
	_, err := strkey.Decode(strkey.VersionByteAccountID, rawAddress)
	if err != nil {
		return "", validation.Errors{
			name: err,
		}
	}
	return rawAddress, nil
}
