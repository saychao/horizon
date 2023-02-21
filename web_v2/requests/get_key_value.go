package requests

import "net/http"

type GetKeyValue struct {
	*base
	Key string
}

func NewGetKeyValue(r *http.Request) (*GetKeyValue, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	key := b.getString("key")

	return &GetKeyValue{
		base: b,
		Key:  key,
	}, nil
}
