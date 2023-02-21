package history2

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/regources/generated"
)

// CreatePoll - represents instance of voting campaign
type Poll struct {
	ID                       int64               `db:"id"`
	PermissionType           uint32              `db:"permission_type"`
	NumberOfChoices          uint32              `db:"number_of_choices"`
	Data                     PollData            `db:"data"`
	StartTime                time.Time           `db:"start_time"`
	EndTime                  time.Time           `db:"end_time"`
	OwnerID                  string              `db:"owner_id"`
	ResultProviderID         string              `db:"result_provider_id"`
	VoteConfirmationRequired bool                `db:"vote_confirmation_required"`
	Details                  regources.Details   `db:"details"`
	State                    regources.PollState `db:"state"`
}

type PollData regources.PollData

//Value - implements db driver method for auto marshal
func (r PollData) Value() (driver.Value, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal poll data")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *PollData) Scan(src interface{}) error {
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
		return errors.Wrap(err, "failed to unmarshal poll data")
	}

	return nil
}
