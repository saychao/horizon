package history2

import "time"

type AssetPair struct {
	Base            string    `db:"base"`
	Quote           string    `db:"quote"`
	CurrentPrice    int64     `db:"current_price"`
	LedgerCloseTime time.Time `db:"ledger_close_time"`
}

func NewAssetPair(base, quote string, currentPrice int64, ledgerCloseTime time.Time) AssetPair {
	return AssetPair{
		Base:            base,
		Quote:           quote,
		CurrentPrice:    currentPrice,
		LedgerCloseTime: ledgerCloseTime,
	}
}
