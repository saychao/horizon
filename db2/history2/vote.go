package history2

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// Vote - represents choice of voting campaign participant
type Vote struct {
	ID       int64    `db:"id"`
	PollID   int64    `db:"poll_id"`
	VoterID  string   `db:"voter_id"`
	VoteData VoteData `db:"data"`
}

type VoteData regources.VoteData

//Value - implements db driver method for auto marshal
func (r VoteData) Value() (driver.Value, error) {
	result, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal poll data")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *VoteData) Scan(src interface{}) error {
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
