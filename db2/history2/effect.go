package history2

import (
	"database/sql/driver"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

// EffectType describe the effect of some operation to the account or particular balance
type EffectType int64

const (
	// EffectTypeFunded - balance received funds from other balance
	EffectTypeFunded EffectType = iota + 1
	// EffectTypeIssued - funds have been issued to the balance
	EffectTypeIssued
	// EffectTypeCharged - balance has been charged
	EffectTypeCharged
	// EffectTypeWithdrawn - balance has been charged from locked and corresponding amount of tokens has been destroyed
	EffectTypeWithdrawn
	// EffectTypeLocked - funds has been locked on the balance
	EffectTypeLocked
	// EffectTypeUnlocked - funds has been unlocked on the balance
	EffectTypeUnlocked
	// EffectTypeChargedFromLocked - funds has been charged from locked amount on balance
	EffectTypeChargedFromLocked
	// EffectTypeMatched - balance has been charged or received funds due to match of the offers
	EffectTypeMatched
)

// Effect stores the details of the operation effect on balance in union switch form. Only one value should be selected
// Effect should never store more than one change to the account or balance
type Effect struct {
	//NOTE: omitempty MUST be specified for each switch value
	Type              EffectType           `json:"type"`
	Funded            *BalanceChangeEffect `json:"funded,omitempty"`
	Issued            *BalanceChangeEffect `json:"issued,omitempty"`
	Charged           *BalanceChangeEffect `json:"charged,omitempty"`
	Withdrawn         *BalanceChangeEffect `json:"withdrawn,omitempty"`
	Locked            *BalanceChangeEffect `json:"locked,omitempty"`
	Unlocked          *BalanceChangeEffect `json:"unlocked,omitempty"`
	ChargedFromLocked *BalanceChangeEffect `json:"charged_from_locked,omitempty"`
	Matched           *MatchEffect         `json:"matched,omitempty"`
}

//Value - converts effect into jsonb
func (r Effect) Value() (driver.Value, error) {
	result, err := pgdb.JSONValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal effect")
	}

	return result, nil
}

//Scan - converts json into Effect
func (r *Effect) Scan(src interface{}) error {
	err := pgdb.JSONScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan effect")
	}

	return nil
}

// MatchEffect - describes changes to base and quote balance occurred on match
type MatchEffect struct {
	OfferID     int64                         `json:"offer_id"`
	OrderBookID int64                         `json:"order_book_id"`
	Price       regources.Amount              `json:"price"`
	Charged     ParticularBalanceChangeEffect `json:"charged"`
	Funded      ParticularBalanceChangeEffect `json:"funded"`
}

// ParticularBalanceChangeEffect - describes movement of fund for particular balance
type ParticularBalanceChangeEffect struct {
	BalanceChangeEffect
	BalanceAddress string `json:"balance_address"`
	AssetCode      string `json:"asset_code"`
}

// BalanceChangeEffect - describes movement of funds
type BalanceChangeEffect struct {
	Amount regources.Amount `json:"amount"`
	Fee    regources.Fee    `json:"fee"`
}
