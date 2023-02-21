package resources

import (
	"github.com/SafeRE-IT/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewDeferredPaymentKey(ID int64) regources.Key {
	return regources.NewKeyInt64(ID, regources.DEFERRED_PAYMENTS)
}

func NewDeferredPayment(record history2.DeferredPayment) regources.DeferredPayment {
	return regources.DeferredPayment{
		Key: NewDeferredPaymentKey(record.ID),
		Attributes: regources.DeferredPaymentAttributes{
			Amount:  record.Amount,
			Details: regources.Details(record.Details),
			StateI:  int32(record.State),
			State:   record.State.String(),
		},
		Relationships: regources.DeferredPaymentRelationships{
			Destination:   NewAccountKey(record.DestinationAccount).AsRelation(),
			Source:        NewAccountKey(record.SourceAccount).AsRelation(),
			SourceBalance: NewBalanceKey(record.SourceBalance).AsRelation(),
		},
	}
}
