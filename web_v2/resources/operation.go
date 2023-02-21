package resources

import (
	"strconv"

	"github.com/SafeRE-IT/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewOperation - creates new instance of regources.Operation
func NewOperation(op history2.Operation) regources.Operation {
	return regources.Operation{
		Key: NewOperationKey(op.ID),
		Attributes: regources.OperationAttributes{
			AppliedAt: op.LedgerCloseTime,
		},
		Relationships: regources.OperationRelationships{
			Tx:      NewTxKey(op.TxID).AsRelation(),
			Source:  NewAccountKey(op.Source).AsRelation(),
			Details: NewOperationDetails(op).GetKey().AsRelation(),
		},
	}
}

// NewOperationKey - creates new key for operation
func NewOperationKey(opID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(opID, 10),
		Type: regources.OPERATIONS,
	}
}
