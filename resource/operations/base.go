package operations

import (
	"fmt"
	"github.com/SafeRE-IT/horizon/httpx"
	"github.com/SafeRE-IT/horizon/render/hal"
	"strconv"

	"time"

	"github.com/SafeRE-IT/horizon/db2/history"
	"github.com/SafeRE-IT/horizon/resource/base"
	"gitlab.com/tokend/go/amount"
	"golang.org/x/net/context"
)

// Base represents the common attributes of an operation resource
type Base struct {
	Links struct {
		Self        hal.Link `json:"self"`
		Transaction hal.Link `json:"transaction"`
		Succeeds    hal.Link `json:"succeeds"`
		Precedes    hal.Link `json:"precedes"`
	} `json:"_links"`

	ID                  string             `json:"id"`
	PT                  string             `json:"paging_token"`
	TransactionID       string             `json:"transaction_id"`
	SourceAccount       string             `json:"source_account,omitempty"`
	Type                string             `json:"type"`
	TypeI               int32              `json:"type_i"`
	StateI              int32              `json:"state_i"`
	State               string             `json:"state"`
	Identifier          string             `json:"identifier"`
	LedgerCloseTime     time.Time          `json:"ledger_close_time"`
	Participants        []base.Participant `json:"participants,omitempty"`
	OperationFee        string             `json:"operation_fee"`
	TransactionFeeAsset string             `json:"operation_fee_asset,omitempty"`
}

// PagingToken implements hal.Pageable
func (this Base) PagingToken() string {
	return this.PT
}

// Populate fills out this resource using `row` as the source.
func (this *Base) Populate(
	ctx context.Context, row history.Operation, participants []*history.Participant, public bool,
) error {
	this.ID = fmt.Sprintf("%d", row.ID)
	this.PT = row.PagingToken()
	this.TransactionID = fmt.Sprintf("%d", row.TransactionID)
	this.SourceAccount = row.SourceAccount
	this.populateType(row)
	this.LedgerCloseTime = row.LedgerCloseTime
	this.Participants = make([]base.Participant, len(participants))
	this.StateI = int32(row.State)
	this.State = row.State.String()
	this.Identifier = strconv.FormatInt(row.Identifier, 10)
	this.OperationFee = amount.String(0)
	for i := range participants {
		err := this.Participants[i].Populate(participants[i], row.Type, public)
		if err != nil {
			return err
		}
	}

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	self := fmt.Sprintf("/operations/%d", row.ID)
	this.Links.Self = lb.Link(self)
	this.Links.Succeeds = lb.Linkf(nil, "/effects?order=desc&cursor=%s", this.PT)
	this.Links.Precedes = lb.Linkf(nil, "/effects?order=asc&cursor=%s", this.PT)
	this.Links.Transaction = lb.Linkf(nil, "/transactions/")

	if public {
		this.SourceAccount = ""
	}
	return nil
}

func (this *Base) populateType(row history.Operation) {
	this.TypeI = int32(row.Type)
	this.Type = row.Type.ShortString()
}
