// Package resource contains the type definitions for all of horizons
// response resources.
package resource

import (
	"fmt"
	"time"

	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/render/hal"
	"github.com/saychao/horizon/resource/operations"
	"golang.org/x/net/context"
)

type BalancePublic struct {
	ID        string `json:"id,omitempty"`
	BalanceID string `json:"balance_id"`
	AccountID string `json:"account_id"`
	Asset     string `json:"asset"`
}

// Balance represents an account's holdings for a single currency type
type Balance struct {
	BalancePublic
	Balance          string `json:"balance,omitempty"`
	Locked           string `json:"locked,omitempty"`
	RequireReview    bool   `json:"require_review"`
	AssetDetails     *Asset `json:"asset_details,omitempty"`
	ConvertedBalance string `json:"converted_balance,omitempty"`
	ConvertedLocked  string `json:"converted_locked,omitempty"`
	ConvertedToAsset string `json:"converted_to_asset,omitempty"`
}

// HistoryAccount is a simple resource, used for the account collection actions.
// It provides only the "TotalOrderID" of the account and its account id.
type HistoryAccount struct {
	ID        string `json:"id,omitempty"`
	PT        string `json:"paging_token,omitempty"`
	AccountID string `json:"account_id"`
}

// TransactionSuccess represents the result of a successful transaction
// submission.
type TransactionSuccess struct {
	Links struct {
		Transaction hal.Link `json:"transaction"`
	} `json:"_links"`
	Hash   string `json:"hash"`
	Ledger int32  `json:"ledger"`
	Env    string `json:"envelope_xdr"`
	Result string `json:"result_xdr"`
	Meta   string `json:"result_meta_xdr"`
}

// NewOperation returns a resource of the appropriate sub-type for the provided
// operation record.
func NewOperation(
	ctx context.Context,
	row history.Operation,
	participants []*history.Participant,
) (result hal.Pageable, err error) {
	return operations.New(ctx, row, participants, false)
}

func NewPublicOperation(
	ctx context.Context, row history.Operation, participants []*history.Participant,
) (result hal.Pageable, err error) {
	return operations.New(ctx, row, participants, true)
}

type Data struct {
	Ledgers []DataLedger `json:"ledgers"`
}

type DataLedger struct {
	ClosedAt     time.Time               `json:"cloased_at"`
	Sequence     int32                   `json:"sequence"`
	LedgerHash   string                  `json:"ledger_hash"`
	Transactions []DataLedgerTransaction `json:"transactions"`
}

func (d DataLedger) PagingToken() string {
	return fmt.Sprintf("%d", d.Sequence)
}

type DataLedgerTransaction struct {
	ID         int64          `json:"id"`
	Operations []hal.Pageable `json:"operations"`
}
