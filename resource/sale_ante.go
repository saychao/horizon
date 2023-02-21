package resource

import (
	"strconv"

	"github.com/SafeRE-IT/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
)

type SaleAnte struct {
	PT                   string `json:"paging_token"`
	SaleID               string `json:"sale_id"`
	ParticipantBalanceID string `json:"participant_balance_id"`
	Amount               string `json:"amount"`
	AssetCode            string `json:"asset_code"`
}

func (s *SaleAnte) Populate(raw *core.SaleAnte, assetCode string) {
	s.SaleID = strconv.FormatUint(raw.SaleID, 10)
	s.ParticipantBalanceID = raw.ParticipantBalanceID
	s.Amount = amount.StringU(raw.Amount)
	s.AssetCode = assetCode
}

func (s *SaleAnte) PagingToken() string {
	return s.PT
}
