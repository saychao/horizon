package resources

import (
	"strconv"

	"gitlab.com/tokend/go/xdr"

	"github.com/saychao/horizon/db2/core2"

	regources "gitlab.com/tokend/regources/generated"
)

// NewAccountRule - returns a new instance of account rule
func NewAccountRule(rule core2.AccountRule) regources.AccountRule {
	return regources.AccountRule{
		Key: NewAccountRuleKey(rule.ID),
		Attributes: regources.AccountRuleAttributes{
			Resource: rule.Resource,
			Action:   xdr.AccountRuleAction(rule.Action),
			Forbids:  rule.Forbids,
			Details:  rule.Details,
		},
	}
}

func NewAccountRuleKey(id uint64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatUint(id, 10),
		Type: regources.ACCOUNT_RULES,
	}
}

// NewSignerRule - returns a new instance of signer rule
func NewSignerRule(rule core2.SignerRule) regources.SignerRule {
	return regources.SignerRule{
		Key: regources.Key{
			ID:   strconv.FormatUint(rule.ID, 10),
			Type: regources.SIGNER_RULES,
		},
		Attributes: regources.SignerRuleAttributes{
			Resource:  rule.Resource,
			Action:    xdr.SignerRuleAction(rule.Action),
			Forbids:   rule.Forbids,
			IsDefault: rule.IsDefault,
			Details:   rule.Details,
		},
		Relationships: regources.SignerRuleRelationships{
			Owner: NewAccountKey(rule.OwnerID).AsRelation(),
		},
	}
}
