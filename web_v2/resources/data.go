package resources

import (
	"github.com/saychao/horizon/db2/history2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewDataKey(id int64) regources.Key {
	return regources.NewKeyInt64(id, regources.DATAS)
}

func NewData(data history2.Data) regources.Data {
	return regources.Data{
		Key: NewDataKey(data.ID),
		Attributes: regources.DataAttributes{
			Type:  uint64(data.Type),
			Value: regources.Details(data.Value),
		},
		Relationships: regources.DataRelationships{
			Owner: NewAccountKey(data.Owner).AsRelation(),
		},
	}
}
