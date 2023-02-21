package resource

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/SafeRE-IT/horizon/db2/core"
	"github.com/SafeRE-IT/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

type Sale struct {
	PT                string                 `json:"paging_token"`
	ID                string                 `json:"id"`
	OwnerID           string                 `json:"owner_id"`
	BaseAsset         string                 `json:"base_asset"`
	DefaultQuoteAsset string                 `json:"default_quote_asset"`
	StartTime         time.Time              `json:"start_time"`
	EndTime           time.Time              `json:"end_time"`
	SoftCap           string                 `json:"soft_cap"`
	HardCap           string                 `json:"hard_cap"`
	Details           map[string]interface{} `json:"details"`
	State             regources.Flag         `json:"state"`
	Statistics        SaleStatistics         `json:"statistics"`
	QuoteAssets       history.QuoteAssets    `json:"quote_assets"`
	BaseHardCap       string                 `json:"base_hard_cap"`
	BaseCurrentCap    string                 `json:"base_current_cap"`
	CurrentCap        string                 `json:"current_cap"`
	SaleType          regources.Flag         `json:"sale_type"`
}

type SaleStatistics struct {
	Investors int `json:"investors"`
}

func (s *Sale) Populate(h *history.Sale) error {
	var details map[string]interface{}
	err := json.Unmarshal(h.Details, &details)

	if err != nil {
		return err
	}

	s.PT = strconv.FormatUint(h.ID, 10)
	s.ID = strconv.FormatUint(h.ID, 10)
	s.OwnerID = h.OwnerID
	s.BaseAsset = h.BaseAsset
	s.DefaultQuoteAsset = h.DefaultQuoteAsset
	s.StartTime = h.StartTime
	s.EndTime = h.EndTime
	s.SoftCap = amount.StringU(h.SoftCap)
	s.HardCap = amount.StringU(h.HardCap)
	s.Details = details
	s.State.Name = h.State.String()
	s.State.Value = int32(h.State)
	s.QuoteAssets = h.QuoteAssets
	s.BaseHardCap = amount.String(h.BaseHardCap)
	s.BaseCurrentCap = amount.String(h.BaseCurrentCap)
	s.CurrentCap = h.CurrentCap
	s.SaleType.Name = h.SaleType.ShortString()
	s.SaleType.Value = int32(h.SaleType)

	return nil
}

func (s *Sale) PopulateStat(offers []core.Offer, balances []core.Balance) error {
	if len(offers) == 0 && len(balances) == 0 {
		return nil
	}

	uniqueInvestors := make(map[string]bool)
	for _, offer := range offers {
		uniqueInvestors[offer.OwnerID] = true
	}

	for _, balance := range balances {
		if balance.Amount == 0 {
			continue
		}
		uniqueInvestors[balance.AccountID] = true
	}

	quantity := len(uniqueInvestors)
	s.Statistics.Investors = quantity
	return nil
}

func divToAmountStr(sum *big.Int, quantity int) string {
	result := sum.Div(sum, big.NewInt(int64(quantity)))
	resultStr := fmt.Sprintf("%d", math.MaxInt64)
	if result.IsInt64() {
		resultStr = amount.String(result.Int64())
	}
	return resultStr
}

// PagingToken implementation for hal.Pageable
func (s *Sale) PagingToken() string {
	return s.PT
}
