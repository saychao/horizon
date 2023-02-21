package history

import (
	"encoding/json"
	"github.com/saychao/horizon/db2"

	"time"

	"github.com/go-errors/errors"
	"github.com/guregu/null"
	"gitlab.com/tokend/go/xdr"
)

// Operation is a row of data from the `history_operations` table
type Operation struct {
	db2.TotalOrderID
	TransactionID    int64             `db:"transaction_id"`
	ApplicationOrder int32             `db:"application_order"`
	Type             xdr.OperationType `db:"type"`
	DetailsString    null.String       `db:"details"`
	LedgerCloseTime  time.Time         `db:"ledger_close_time"`
	SourceAccount    string            `db:"source_account"`
	State            OperationState    `db:"state"`
	Identifier       int64             `db:"identifier"`
}

// UnmarshalDetails unmarshals the details of this operation into `dest`
func (r *Operation) UnmarshalDetails(dest interface{}) error {
	if !r.DetailsString.Valid {
		return nil
	}

	err := json.Unmarshal([]byte(r.DetailsString.String), &dest)
	if err != nil {
		err = errors.Wrap(err, 1)
	}

	return err
}
