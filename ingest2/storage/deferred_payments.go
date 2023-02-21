package storage

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

// DeferredPayment is helper struct to operate with `deferredPayments`
type DeferredPayment struct {
	repo    *pgdb.DB
	q       history2.DeferredPaymentQ
	updater sq.UpdateBuilder
}

// NewDeferredPayment - creates new instance of the `CreateDeferredPayment`
func NewDeferredPayment(repo *pgdb.DB) *DeferredPayment {
	return &DeferredPayment{
		repo:    repo,
		q:       history2.NewDeferredPaymentQ(repo),
		updater: sq.Update("deferred_payments"),
	}
}

// Insert - inserts new deferredPayment
func (q *DeferredPayment) Insert(deferredPayment history2.DeferredPayment) error {
	sql := sq.Insert("deferred_payments").
		Columns(
			"id",
			"amount",
			"details",
			"source_account",
			"source_balance",
			"destination_account",
			"state",
		).
		Values(
			deferredPayment.ID,
			deferredPayment.Amount,
			deferredPayment.Details,
			deferredPayment.SourceAccount,
			deferredPayment.SourceBalance,
			deferredPayment.DestinationAccount,
			deferredPayment.State,
		)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert deferredPayment", logan.F{"deferredPayment_id": deferredPayment.ID})
	}

	return nil
}

// Update - updates existing deferredPayment
func (q *DeferredPayment) Update(deferredPayment history2.DeferredPayment) error {
	sql := sq.Update("deferred_payments").SetMap(map[string]interface{}{
		"amount":              deferredPayment.Amount,
		"source_account":      deferredPayment.SourceAccount,
		"source_balance":      deferredPayment.SourceBalance,
		"destination_account": deferredPayment.DestinationAccount,
	}).Where(sq.Eq{"id": deferredPayment.ID})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update deferredPayment", logan.F{"deferredPayment_id": deferredPayment.ID})
	}

	return nil
}

func (q *DeferredPayment) UpdateState(id int64, state history2.DeferredPaymentState) error {
	sql := sq.Update("deferred_payments").
		Set("state", state).
		Where(sq.Eq{"id": id})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update deferredPayment", logan.F{"deferredPayment_id": id})
	}

	return nil
}

func (q *DeferredPayment) MustDeferredPayment(id int64) history2.DeferredPayment {
	deferredPayment, err := q.getDeferredPayment(id)
	if err != nil {
		panic(err)
	}

	return deferredPayment
}

func (q *DeferredPayment) getDeferredPayment(id int64) (history2.DeferredPayment, error) {
	deferredPayment, err := q.q.GetByID(id)
	if err != nil {
		return history2.DeferredPayment{}, errors.Wrap(err, "failed to get deferredPayment by id", logan.F{
			"deferredPayment_id": id,
		})
	}

	if deferredPayment == nil {
		return history2.DeferredPayment{}, errors.From(errors.New("deferredPayment missing"), logan.F{
			"deferredPayment_id": id,
		})
	}

	return *deferredPayment, nil
}
