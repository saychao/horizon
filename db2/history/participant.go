package history

import (
	"encoding/json"

	"github.com/go-errors/errors"
	"gitlab.com/tokend/go/xdr"
)

// helper struct, should exists mostly in operation_q
type Participant struct {
	OperationID int64  `db:"history_operation_id"`
	AccountID   string `db:"account_id"`
	BalanceID   string `db:"balance_id"`
	Nickname    string
	Email       string
	Mobile      string
	Details     []byte
	UserType    string
	Effects     *[]byte
}

type OperationParticipants struct {
	OpType       xdr.OperationType
	Participants []*Participant
}

func (p *Participant) UnmarshalEffects(dest interface{}) error {
	err := json.Unmarshal((*p.Effects), &dest)
	if err != nil {
		err = errors.Wrap(err, 1)
	}

	return err
}
