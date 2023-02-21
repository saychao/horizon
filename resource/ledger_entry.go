package resource

import "gitlab.com/tokend/go/xdr"

type LedgerEntry struct {
	LastModifiedLedgerSeq uint32        `json:"last_modified_ledger_seq"`
	TypeI                 int32         `json:"type_i"`
	Type                  string        `json:"type"`
	Account               *AccountEntry `json:"account"`
	Asset                 *AssetEntry   `json:"asset"`
	Balance               *BalanceEntry `json:"balance"`
}

func (r *LedgerEntry) Populate(entry xdr.LedgerEntry) bool {
	r.TypeI = int32(entry.Data.Type)
	r.Type = entry.Data.Type.ShortString()
	r.LastModifiedLedgerSeq = uint32(entry.LastModifiedLedgerSeq)

	switch entry.Data.Type {
	case xdr.LedgerEntryTypeAccount:
		r.Account = new(AccountEntry)
		r.Account.Populate(entry.Data.MustAccount())
	case xdr.LedgerEntryTypeAsset:
		r.Asset = new(AssetEntry)
		r.Asset.Populate(entry.Data.MustAsset())
	case xdr.LedgerEntryTypeBalance:
		r.Balance = new(BalanceEntry)
		r.Balance.Populate(entry.Data.MustBalance())
	default:
		return false
	}
	return true
}

type LedgerKey struct {
	TypeI   int32             `json:"type_i"`
	Type    string            `json:"type"`
	Account *LedgerKeyAccount `json:"account"`
	Asset   *LedgerKeyAsset   `json:"asset"`
	Balance *LedgerKeyBalance `json:"balance"`
}

func (r *LedgerKey) Populate(ledgerKey xdr.LedgerKey) bool {
	r.TypeI = int32(ledgerKey.Type)
	r.Type = ledgerKey.Type.ShortString()

	switch ledgerKey.Type {
	case xdr.LedgerEntryTypeAccount:
		r.Account = new(LedgerKeyAccount)
		r.Account.Populate(ledgerKey.MustAccount())
	case xdr.LedgerEntryTypeAsset:
		r.Asset = new(LedgerKeyAsset)
		r.Asset.Populate(ledgerKey.MustAsset())
	case xdr.LedgerEntryTypeBalance:
		r.Balance = new(LedgerKeyBalance)
		r.Balance.Populate(ledgerKey.MustBalance())
	default:
		return false
	}
	return true
}
