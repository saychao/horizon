package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// CreateSwap is helper struct to operate with `swaps`
type Swap struct {
	repo *pgdb.DB

	swapQ history2.SwapsQ
}

// NewSwap - creates new instance of the `CreateSwap`
func NewSwap(repo *pgdb.DB) *Swap {
	return &Swap{
		repo:  repo,
		swapQ: history2.NewSwapsQ(repo),
	}
}

// Insert - inserts new swap
func (q *Swap) Insert(swap history2.Swap) error {

	sql := sq.Insert("swaps").
		SetMap(map[string]interface{}{
			"id":                      swap.ID,
			"source_account":          swap.SourceAccount,
			"source_balance":          swap.SourceBalance,
			"destination_account":     swap.DestinationAccount,
			"destination_balance":     swap.DestinationBalance,
			"asset":                   swap.Asset,
			"amount":                  swap.Amount,
			"created_at":              swap.CreatedAt,
			"lock_time":               swap.LockTime,
			"secret_hash":             swap.SecretHash,
			"secret":                  swap.Secret,
			"source_fixed_fee":        swap.SourceFixedFee,
			"source_percent_fee":      swap.SourcePercentFee,
			"destination_fixed_fee":   swap.DestinationFixedFee,
			"destination_percent_fee": swap.DestinationPercentFee,
			"details":                 swap.Details,
			"state":                   swap.State,
		})

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert swap", logan.F{"swap_id": swap.ID})
	}

	return nil
}

// SetState - sets state
func (q *Swap) SetState(id int64, state regources.SwapState) error {
	sql := sq.Update("swaps").Set("state", state).Where("id = ?", id)
	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"swap_id": id})
	}

	return nil
}

// SetSecret - sets secret
func (q *Swap) SetSecret(id int64, secret string) error {
	sql := sq.Update("swaps").Set("secret", secret).Where("id = ?", id)
	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set secret", logan.F{"swap_id": id})
	}

	return nil
}

func (q *Swap) getSwap(id int64) (history2.Swap, error) {
	swap, err := q.swapQ.GetByID(id)
	if err != nil {
		return history2.Swap{}, errors.Wrap(err, "failed to get swap by id", logan.F{
			"swap_id": id,
		})
	}

	if swap == nil {
		return history2.Swap{}, errors.From(errors.New("swap missing"), logan.F{
			"swap_id": id,
		})
	}

	return *swap, nil
}

func (q *Swap) MustSwap(id int64) history2.Swap {
	swap, err := q.getSwap(id)
	if err != nil {
		panic(err)
	}

	return swap
}
