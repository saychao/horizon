package resource

import "github.com/SafeRE-IT/horizon/db2/history"

type PriceHistory struct {
	Prices []history.PricePoint `json:"prices"`
}
