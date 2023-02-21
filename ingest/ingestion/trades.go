package ingestion

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

func (ingest *Ingestion) StoreTrades(orderBookID uint64, result xdr.ManageOfferSuccessResult, ledgerCloseTime int64) error {
	if len(result.OffersClaimed) == 0 {
		return nil
	}

	q := ingest.trades
	for i := range result.OffersClaimed {
		claimed := result.OffersClaimed[i]
		q = q.Values(orderBookID, string(result.BaseAsset),
			string(result.QuoteAsset), int64(claimed.BaseAmount),
			int64(claimed.QuoteAmount), int64(claimed.CurrentPrice), time.Unix(ledgerCloseTime, 0).UTC())
	}

	err := ingest.DB.Exec(q)
	if err != nil {
		return errors.Wrap(err, "failed to insert into trades")
	}

	return nil
}
