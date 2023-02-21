package ingest

import (
	"encoding/json"
	regources "gitlab.com/tokend/regources/generated"
	"time"

	"github.com/saychao/horizon/db2/history"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
)

func saleCreate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	rawSale := ledgerEntry.Data.MustSale()
	sale, err := convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale")
	}

	// if sale already exists - it was in state "PROMOTION" and we need to update it
	histSale, err := is.Ingestion.HistoryQ().Sales().ByID(sale.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get sale from History DB")
	}

	if histSale != nil {
		err = is.Ingestion.HistoryQ().Sales().Update(*sale)
		if err != nil {
			return errors.Wrap(err, "failed to update sale")
		}
		return nil
	}

	err = is.Ingestion.HistoryQ().Sales().Insert(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale")
	}

	return nil
}

func saleUpdate(is *Session, ledgerEntry *xdr.LedgerEntry) error {
	rawSale := ledgerEntry.Data.MustSale()
	sale, err := convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale")
	}

	err = is.Ingestion.HistoryQ().Sales().Update(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to update sale")
	}

	return nil
}

func convertSale(raw xdr.SaleEntry) (*history.Sale, error) {
	var quoteAssets []history.QuoteAsset
	for i := range raw.QuoteAssets {
		quoteAssets = append(quoteAssets, history.QuoteAsset{
			Asset:          string(raw.QuoteAssets[i].QuoteAsset),
			Price:          amount.StringU(uint64(raw.QuoteAssets[i].Price)),
			QuoteBalanceID: raw.QuoteAssets[i].QuoteBalance.AsString(),
			CurrentCap:     amount.StringU(uint64(raw.QuoteAssets[i].CurrentCap)),
		})
	}

	var saleDetails regources.Details
	_ = json.Unmarshal([]byte(raw.Details), &saleDetails)

	saleType := raw.SaleTypeExt.SaleType

	return &history.Sale{
		ID:                uint64(raw.SaleId),
		OwnerID:           raw.OwnerId.Address(),
		BaseAsset:         string(raw.BaseAsset),
		DefaultQuoteAsset: string(raw.DefaultQuoteAsset),
		StartTime:         time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:           time.Unix(int64(raw.EndTime), 0).UTC(),
		SoftCap:           uint64(raw.SoftCap),
		HardCap:           uint64(raw.HardCap),
		Details:           saleDetails,
		State:             history.SaleStateOpen,
		QuoteAssets: history.QuoteAssets{
			QuoteAssets: quoteAssets,
		},
		BaseCurrentCap: int64(raw.CurrentCapInBase),
		BaseHardCap:    int64(raw.MaxAmountToBeSold),
		SaleType:       saleType,
	}, nil
}

func (is *Session) processCancelSaleCreationRequest(
	op xdr.CancelSaleCreationRequestOp,
	result xdr.CancelSaleCreationRequestResult,
) error {
	if result.Code != xdr.CancelSaleCreationRequestResultCodeSuccess {
		return nil
	}

	err := is.Ingestion.HistoryQ().ReviewableRequests().
		Cancel(uint64(op.RequestId))
	if err != nil {
		return errors.Wrap(err,
			"failed to update sale creation request state to cancel", logan.F{
				"request_id": uint64(op.RequestId),
			})
	}

	return nil
}
