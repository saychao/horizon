package history

import (
	"github.com/SafeRE-IT/horizon/db2"
	regources "gitlab.com/tokend/regources/generated"
	"time"

	"github.com/lib/pq"
)

type Contract struct {
	db2.TotalOrderID
	Contractor      string            `db:"contractor"`
	Customer        string            `db:"customer"`
	Escrow          string            `db:"escrow"`
	StartTime       time.Time         `db:"start_time"`
	EndTime         time.Time         `db:"end_time"`
	InitialDetails  regources.Details `db:"initial_details"`
	CustomerDetails regources.Details `db:"customer_details"`
	Invoices        pq.Int64Array     `db:"invoices"`
	State           int32             `db:"state"`
}

type ContractDetails struct {
	ContractID int64             `db:"contract_id"`
	Details    regources.Details `db:"details"`
	Author     string            `db:"author"`
	CreatedAt  time.Time         `db:"created_at"`
}

type ContractDispute struct {
	ContractID int64             `db:"contract_id"`
	Reason     regources.Details `db:"reason"`
	Author     string            `db:"author"`
	CreatedAt  time.Time         `db:"created_at"`
}
