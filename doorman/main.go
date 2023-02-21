package doorman

import (
	"github.com/saychao/horizon/db2/core2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/resources"
)

// SignersQ - helpers struct to fetch signers from storage. Implements doorman interface
type SignersQ struct {
	q core2.SignerQ
}

// NewSignersQ - creates new instance of SignersQ
func NewSignersQ(q core2.SignerQ) *SignersQ {
	return &SignersQ{
		q: q,
	}
}

// Signers get account signers, nil slice is returned if account is not found
func (q *SignersQ) Signers(address string) ([]resources.Signer, error) {
	signers, err := q.q.FilterByAccountAddress(address).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load signers for address", logan.F{
			"address": address,
		})
	}

	result := make([]resources.Signer, 0, len(signers))
	for _, signer := range signers {
		result = append(result, resources.Signer{
			// FIXME with it properly in doorman as well.
			AccountID: signer.PublicKey,
			Weight:    int(signer.Weight),
			Role:      signer.RoleID,
			Identity:  signer.Identity,
		})
	}

	return result, nil
}
