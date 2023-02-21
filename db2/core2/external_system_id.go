package core2

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	regources "gitlab.com/tokend/regources/generated"
)

type ExternalSystemID struct {
	ID                 uint64             `db:"id"`
	AccountID          string             `db:"account_id"`
	ExternalSystemType int32              `db:"external_system_type"`
	Data               externalSystemData `db:"data"`
	IsDeleted          bool               `db:"is_deleted"`
	ExpiresAt          time.Time          `db:"expires_at"`
	BindedAt           time.Time          `db:"binded_at"`
}

type externalSystemData regources.ExternalSystemData

//Value - implements db driver method for auto marshal
func (r externalSystemData) Value() (driver.Value, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal ExternalSystemData data")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *externalSystemData) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	default:
		return errors.New("Unexpected type for jsonb")
	}

	err := json.Unmarshal(data, r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ExternalSystemData data")
	}

	return nil
}
