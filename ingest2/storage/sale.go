package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/SafeRE-IT/horizon/db2/history2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// Sale is helper struct to operate with `sales`
type Sale struct {
	repo *pgdb.DB
}

// NewSale - creates new instance of the `Sale`
func NewSale(repo *pgdb.DB) *Sale {
	return &Sale{
		repo: repo,
	}
}

// Insert - inserts new sale
func (q *Sale) Insert(sale history2.Sale) error {
	sql := sq.Insert("sales").
		Columns(
			"id", "owner_address", "base_asset", "default_quote_asset", "start_time", "end_time",
			"quote_assets", "soft_cap", "hard_cap", "details", "base_current_cap",
			"base_hard_cap", "sale_type", "state", "version", "access_definition_type",
		).
		Values(
			sale.ID, sale.OwnerAddress, sale.BaseAsset, sale.DefaultQuoteAsset, sale.StartTime, sale.EndTime,
			sale.QuoteAssets, sale.SoftCap, sale.HardCap, sale.Details,
			sale.BaseCurrentCap, sale.BaseHardCap, sale.SaleType, sale.State, sale.Version, sale.AccessDefinitionType,
		)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale", logan.F{"sale_id": sale.ID})
	}

	return nil
}

// Update - updates existing sale
func (q *Sale) Update(sale history2.Sale) error {
	sql := sq.Update("sales").SetMap(map[string]interface{}{
		"owner_address":       sale.OwnerAddress,
		"base_asset":          sale.BaseAsset,
		"default_quote_asset": sale.DefaultQuoteAsset,
		"start_time":          sale.StartTime,
		"end_time":            sale.EndTime,
		"quote_assets":        sale.QuoteAssets,
		"soft_cap":            sale.SoftCap,
		"hard_cap":            sale.HardCap,
		"details":             sale.Details,
		"base_hard_cap":       sale.BaseHardCap,
		"base_current_cap":    sale.BaseCurrentCap,
		"sale_type":           sale.SaleType,
		"state":               sale.State,
		"version":             sale.Version,
	}).Where("id = ?", sale.ID)

	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update sale", logan.F{"sale_id": sale.ID})
	}

	return nil
}

// SetState - sets state
func (q *Sale) SetState(saleID uint64, state regources.SaleState) error {
	sql := sq.Update("sales").Set("state", state).Where("id = ?", saleID)
	err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"sale_id": saleID})
	}

	return nil
}
