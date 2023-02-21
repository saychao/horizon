package resource

import (
	"github.com/SafeRE-IT/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

type Policy struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

// Signer represents one of an account's signers.
type Policies struct {
	PolicyI  int32    `json:"policy"`
	Policies []Policy `json:"policies"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (p *Policies) Populate(row core.Asset) {
	p.PopulateFromInt32(row.Policies)
}

func (p *Policies) PopulateFromInt32(policies int32) {
	p.PolicyI = policies
	for _, policy := range xdr.AssetPolicyAll {
		if (int32(policy) & p.PolicyI) != 0 {
			p.Policies = append(p.Policies, Policy{
				Value: int32(policy),
				Name:  policy.String(),
			})
		}
	}
}
