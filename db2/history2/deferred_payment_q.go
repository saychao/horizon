package history2

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// DeferredPaymentQ is a helper struct to aid in configuring queries that loads deferredPayments
type DeferredPaymentQ struct {
	repo     *pgdb.DB
	selector sq.SelectBuilder
}

// NewDeferredPaymentQ- creates new instance of DeferredPaymentQ
func NewDeferredPaymentQ(repo *pgdb.DB) DeferredPaymentQ {
	return DeferredPaymentQ{
		repo: repo,
		selector: sq.Select(
			"deferred_payments.id",
			"deferred_payments.amount",
			"deferred_payments.details",
			"deferred_payments.source_account",
			"deferred_payments.source_balance",
			"deferred_payments.destination_account",
			"deferred_payments.state",
		).From("deferred_payments deferred_payments"),
	}
}

// GetByID - get deferredPayment by ID
func (q DeferredPaymentQ) GetByID(id int64) (*DeferredPayment, error) {
	q.selector = q.selector.Where(sq.Eq{"deferred_payments.id": id})
	return q.Get()
}

//FilterByDestinationAccount - filters deferredPayments by destination account
func (q DeferredPaymentQ) FilterByDestinationAccount(address string) DeferredPaymentQ {
	q.selector = q.selector.Where(sq.Eq{"deferred_payments.destination_account": address})
	return q
}

//FilterBySourceAccount - filters deferredPayments by source account
func (q DeferredPaymentQ) FilterBySourceAccount(address string) DeferredPaymentQ {
	q.selector = q.selector.Where(sq.Eq{"deferred_payments.source_account": address})
	return q
}

//FilterBySourceBalance - filters deferredPayments by source balance
func (q DeferredPaymentQ) FilterBySourceBalance(address string) DeferredPaymentQ {
	q.selector = q.selector.Where(sq.Eq{"deferred_payments.source_balance": address})
	return q
}

//FilterByAsset - filters deferredPayments by asset code
func (q DeferredPaymentQ) FilterByAsset(asset string) DeferredPaymentQ {
	q.selector = q.selector.Join("balances balances ON balances.id = deferred_payments.source_balance").
		Where(sq.Eq{"balances.asset_code": asset})
	return q
}

//Get - selects deferredPayment from db, returns nil, nil if one does not exists
func (q DeferredPaymentQ) Get() (*DeferredPayment, error) {
	var result DeferredPayment
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load deferredPayment")
	}

	return &result, nil
}

func (q DeferredPaymentQ) Page(params *pgdb.CursorPageParams) DeferredPaymentQ {
	q.selector = params.ApplyTo(q.selector, "deferred_payments.id")
	return q
}

//Select - selects slice from the db, if no deferredPayment found - returns nil, nil
func (q DeferredPaymentQ) Select() ([]DeferredPayment, error) {
	var result []DeferredPayment
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load deferredPayment")
	}

	return result, nil
}
