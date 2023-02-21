package core2

// LedgerHeader is row of data from the `ledgerheaders` table
type License struct {
	ID              int64  `db:"id"`
	Hash            string `db:"hash"`
	PrevLicenseHash string `db:"prev_hash"`
	LedgerHash      string `db:"ledger_hash"`
	AdminCount      int64  `db:"admin_count"`
	DueDate         int64  `db:"due_date"`
}
