package history2

import (
	"database/sql/driver"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

// Sale - represents instance of compounding campaign
type Sale struct {
	ID                   uint64           `db:"id"`
	SoftCap              regources.Amount `db:"soft_cap"`
	HardCap              regources.Amount `db:"hard_cap"`
	BaseCurrentCap       regources.Amount `db:"base_current_cap"`
	BaseHardCap          regources.Amount `db:"base_hard_cap"`
	SaleType             xdr.SaleType     `db:"sale_type"`
	OwnerAddress         string           `db:"owner_address"`
	BaseAsset            string           `db:"base_asset"`
	DefaultQuoteAsset    string           `db:"default_quote_asset"`
	StartTime            time.Time        `db:"start_time"`
	EndTime              time.Time        `db:"end_time"`
	CurrentCap           regources.Amount
	Details              regources.Details                  `db:"details"`
	QuoteAssets          SaleQuoteAssets                    `db:"quote_assets"`
	State                regources.SaleState                `db:"state"`
	AccessDefinitionType regources.SaleAccessDefinitionType `db:"access_definition_type"`
	Version              int32                              `db:"version"`

	*Asset `db:"asset"`
}

//SaleQuoteAssets - assets allowed to invest in sale
type SaleQuoteAssets struct {
	QuoteAssets []SaleQuoteAsset `json:"quote_assets"` // fixme?: move to Sale struct to avoid sale.QuoteAssets.QuoteAssets
}

//Value - implements db driver method for auto marshal
func (r SaleQuoteAssets) Value() (driver.Value, error) {
	result, err := pgdb.JSONValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal quote assets")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *SaleQuoteAssets) Scan(src interface{}) error {
	err := pgdb.JSONScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan quote assets")
	}

	return nil
}

//SaleQuoteAsset - asset allowed to invest into sale
type SaleQuoteAsset struct {
	Asset           string           `json:"asset"`
	Price           regources.Amount `json:"price"`
	QuoteBalanceID  string           `json:"quote_balance_id"`
	CurrentCap      regources.Amount `json:"current_cap"`
	TotalCurrentCap regources.Amount `json:"total_current_cap,omitempty"`
	HardCap         regources.Amount `json:"hard_cap,omitempty"`
}
