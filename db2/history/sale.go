package history

import (
	"database/sql/driver"
	"gitlab.com/distributed_lab/kit/pgdb"
	regources "gitlab.com/tokend/regources/generated"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type Sale struct {
	ID                uint64    `db:"id"`
	OwnerID           string    `db:"owner_id"`
	BaseAsset         string    `db:"base_asset"`
	DefaultQuoteAsset string    `db:"default_quote_asset"`
	StartTime         time.Time `db:"start_time"`
	EndTime           time.Time `db:"end_time"`
	SoftCap           uint64    `db:"soft_cap"`
	HardCap           uint64    `db:"hard_cap"`
	CurrentCap        string
	Details           regources.Details `db:"details"`
	State             SaleState         `db:"state"`
	QuoteAssets       QuoteAssets       `db:"quote_assets"`
	BaseCurrentCap    int64             `db:"base_current_cap"`
	BaseHardCap       int64             `db:"base_hard_cap"`
	SaleType          xdr.SaleType      `db:"sale_type"`
}

type QuoteAssets struct {
	QuoteAssets []QuoteAsset `json:"quote_assets"`
}

func (r QuoteAssets) Value() (driver.Value, error) {
	result, err := pgdb.JSONValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal quote assets")
	}

	return result, nil
}

func (r *QuoteAssets) Scan(src interface{}) error {
	err := pgdb.JSONScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan quote assets")
	}

	return nil
}

type QuoteAsset struct {
	Asset           string `json:"asset"`
	Price           string `json:"price"`
	QuoteBalanceID  string `json:"quote_balance_id"`
	CurrentCap      string `json:"current_cap"`
	TotalCurrentCap string `json:"total_current_cap,omitempty"`
	HardCap         string `json:"hard_cap,omitempty"`
}
