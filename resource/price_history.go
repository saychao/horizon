package resource

import "github.com/saychao/horizon/db2/history"

type PriceHistory struct {
	Prices []history.PricePoint `json:"prices"`
}
