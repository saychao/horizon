package resource

import (
	"fmt"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/history"
)

type LedgerChanges struct {
	ID              string              `json:"id"`
	PT              string              `json:"paging_token"`
	Ledger          int32               `json:"ledger"`
	LedgerCloseTime time.Time           `json:"created_at"`
	Changes         []LedgerEntryChange `json:"changes"`
}

func (lc LedgerChanges) PagingToken() string {
	return lc.PT
}

func (lc *LedgerChanges) Populate(tm history.Transaction) error {
	lc.ID = fmt.Sprintf("%d", tm.ID)
	lc.Ledger = tm.LedgerSequence
	lc.PT = tm.PagingToken()
	lc.LedgerCloseTime = tm.LedgerCloseTime

	txMeta := xdr.TransactionMeta{}
	err := xdr.SafeUnmarshalBase64(tm.TxMeta, &txMeta)
	if err != nil {
		return errors.Wrap(err,
			"failed to unmarshal tx_meta",
			logan.F{
				"id":                tm.ID,
				"ledger_close_time": tm.LedgerCloseTime,
			})
	}

	for _, opMeta := range txMeta.MustOperations() {
		for _, xdrChange := range opMeta.Changes {
			res := LedgerEntryChange{}
			if res.Populate(xdrChange) {
				lc.Changes = append(lc.Changes, res)
			}
		}
	}

	return nil
}

type LedgerEntryChange struct {
	TypeI   int32        `json:"type_i"`
	Type    string       `json:"type"`
	Created *LedgerEntry `json:"created"`
	Updated *LedgerEntry `json:"updated"`
	Removed *LedgerKey   `json:"removed"`
	State   *LedgerEntry `json:"state"`
}

func (r *LedgerEntryChange) Populate(xdrChange xdr.LedgerEntryChange) bool {
	r.TypeI = int32(xdrChange.Type)
	r.Type = xdrChange.Type.ShortString()

	var ok bool
	switch xdrChange.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		r.Created = &LedgerEntry{}
		ok = r.Created.Populate(xdrChange.MustCreated())
	case xdr.LedgerEntryChangeTypeUpdated:
		r.Updated = &LedgerEntry{}
		ok = r.Updated.Populate(xdrChange.MustUpdated())
	case xdr.LedgerEntryChangeTypeRemoved:
		r.Removed = &LedgerKey{}
		ok = r.Removed.Populate(xdrChange.MustRemoved())
	case xdr.LedgerEntryChangeTypeState:
		r.State = &LedgerEntry{}
		ok = r.State.Populate(xdrChange.MustState())
	default:
		return false
	}
	return ok
}
