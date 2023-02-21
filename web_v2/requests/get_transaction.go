package requests

import (
	"net/http"
)

// GetTransaction - represents params to be specified for GetTransaction handler
type GetTransaction struct {
	*base
	// it could be id or hash of the transaction
	ID string
}

// NewGetTransaction returns the new instance of GetTransaction request
func NewGetTransaction(r *http.Request) (*GetTransaction, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeTransactionListLedgerEntryChanges: {},
		},
	})
	if err != nil {
		return nil, err
	}

	return &GetTransaction{
		base: b,
		ID:   b.getString("id"),
	}, nil
}
