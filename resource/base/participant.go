package base

import (
	"fmt"

	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/go/xdr"
)

type BaseEffects interface {
}

type ManageOfferEffects struct {
	BaseAssetCode  string `json:"base_asset_code"`
	QuoteAssetCode string `json:"quote_asset_code"`
	IsBuy          bool   `json:"is_buy"`
	BaseAmount     string `json:"base_amount"`
	QuoteAmount    string `json:"quote_amount"`
	Price          string `json:"price"`
}

type Participant struct {
	AccountID string      `json:"account_id,omitempty"`
	BalanceID string      `json:"balance_id,omitempty"`
	Email     string      `json:"email,omitempty"`
	FullName  string      `json:"full_name,omitempty"`
	Nickname  string      `json:"nickname,omitempty"`
	Effects   BaseEffects `json:"effects,omitempty"`
}

func (f *Participant) Populate(p *history.Participant, opType xdr.OperationType, public bool) error {
	if !public {
		f.AccountID = p.AccountID
		f.BalanceID = p.BalanceID
		f.Nickname = p.Nickname
		f.Email = p.Email
	}
	if p.Effects != nil {
		err := f.PopulateEffects(p, opType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Participant) PopulateEffects(p *history.Participant, opType xdr.OperationType) error {
	var err error
	switch opType {
	case xdr.OperationTypeManageOffer:
		f.Effects = MatchEffects{}
		err = p.UnmarshalEffects(&f.Effects)
	case xdr.OperationTypeCheckSaleState:
		f.Effects = MatchEffects{}
		err = p.UnmarshalEffects(&f.Effects)
	default:
		err = fmt.Errorf("unexpected effects type for op %d: %v", p.OperationID, opType)
	}

	return err
}
