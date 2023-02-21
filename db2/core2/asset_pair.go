package core2

// AssetPair - db representation of asset pair
type AssetPair struct {
	Base                    string `db:"base"`
	Quote                   string `db:"quote"`
	CurrentPrice            int64  `db:"current_price"`
	PhysicalPrice           int64  `db:"physical_price"`
	PhysicalPriceCorrection int64  `db:"physical_price_correction"`
	MaxPriceStep            int64  `db:"max_price_step"`
	Policies                int32  `db:"policies"`

	BaseAsset  *Asset `db:"base_assets"`
	QuoteAsset *Asset `db:"quote_assets"`
}
