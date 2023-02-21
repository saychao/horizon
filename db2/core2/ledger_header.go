package core2

import "gitlab.com/tokend/go/xdr"

// LedgerHeader is row of data from the `ledgerheaders` table
type LedgerHeader struct {
	LedgerHash     string           `db:"ledgerhash"`
	PrevHash       string           `db:"prevhash"`
	BucketListHash string           `db:"bucketlisthash"`
	CloseTime      int64            `db:"closetime"`
	Sequence       int32            `db:"ledgerseq"`
	Version        uint64           `db:"version"`
	Data           xdr.LedgerHeader `db:"data"`
}
