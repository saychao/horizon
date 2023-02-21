package resource

import (
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateTransaction(row history.Transaction) (res regources.Transaction) {
	res.ID = row.TransactionHash
	res.PT = row.PagingToken()
	res.Hash = row.TransactionHash
	res.Ledger = row.LedgerSequence
	res.LedgerCloseTime = row.LedgerCloseTime
	res.Account = row.Account
	res.FeePaid = row.FeePaid
	res.OperationCount = row.OperationCount
	res.EnvelopeXdr = row.TxEnvelope
	res.ResultXdr = row.TxResult
	res.ResultMetaXdr = row.TxMeta
	res.FeeMetaXdr = row.TxFeeMeta
	res.MemoType = row.MemoType
	res.Memo = row.Memo.String
	res.Signatures = strings.Split(row.SignatureString, ",")
	res.ValidBefore = timeString(row.ValidBefore)
	res.ValidAfter = timeString(row.ValidAfter)
	return res
}

func timeString(in null.Int) string {
	if !in.Valid {
		return ""
	}

	return time.Unix(in.Int64, 0).UTC().Format(time.RFC3339)
}
