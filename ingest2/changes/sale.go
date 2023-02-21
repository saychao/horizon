package changes

import (
	"time"

	regources "gitlab.com/tokend/regources/generated"

	history "github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type saleStorage interface {
	//Inserts sale into DB
	Insert(sale history.Sale) error
	//Updates sale
	Update(sale history.Sale) error
	// SetState - sets state
	SetState(saleID uint64, state regources.SaleState) error
}

type saleHandler struct {
	storage      saleStorage
	specificRule accountSpecificRuleStorage
}

func newSaleHandler(storage saleStorage, specificRule accountSpecificRuleStorage) *saleHandler {
	return &saleHandler{
		storage:      storage,
		specificRule: specificRule,
	}
}

// Created - handles creation of new sale
func (c *saleHandler) Created(lc ledgerChange) error {
	rawSale := lc.LedgerChange.MustCreated().Data.MustSale()

	sale, err := c.convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale", logan.F{
			"raw_sale":        rawSale,
			"ledger_sequence": lc.LedgerSeq,
		})
	}
	sale.AccessDefinitionType = c.getAccessDefinitionType(lc)

	err = c.storage.Insert(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale into DB", logan.F{
			"sale": sale,
		})
	}

	return nil
}

// Removed - handles state of the sale due to it was removed
func (c *saleHandler) Removed(lc ledgerChange) error {
	saleID := uint64(lc.LedgerChange.MustRemoved().MustSale().SaleId)
	// sale can be removed by check sale state or cancel sale
	// so we can handle approve of the sale and by default mark is as cancelled
	saleState := regources.SaleStateCanceled
	if lc.OperationResult.Type == xdr.OperationTypeCheckSaleState {
		opEffect := lc.OperationResult.MustCheckSaleStateResult().MustSuccess().Effect
		if opEffect.Effect == xdr.CheckSaleStateEffectClosed {
			saleState = regources.SaleStateClosed
		}
	}

	err := c.storage.SetState(saleID, saleState)
	if err != nil {
		return errors.Wrap(err, "failed to set sale state")
	}

	return nil
}

// Updated - handles update of the sale
func (c *saleHandler) Updated(lc ledgerChange) error {
	rawSale := lc.LedgerChange.MustUpdated().Data.MustSale()
	sale, err := c.convertSale(rawSale)
	if err != nil {
		return errors.Wrap(err, "failed to convert sale ", logan.F{
			"raw_sale":        rawSale,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Update(*sale)
	if err != nil {
		return errors.Wrap(err, "failed to update sale", logan.F{
			"sale": sale,
		})
	}
	return nil
}

func (c *saleHandler) getAccessDefinitionType(lc ledgerChange) regources.SaleAccessDefinitionType {
	for _, change := range lc.OpChanges {
		if change.Type != xdr.LedgerEntryChangeTypeCreated {
			continue
		}
		if change.MustCreated().Data.Type != xdr.LedgerEntryTypeAccountSpecificRule {
			continue
		}

		rule := change.MustCreated().Data.MustAccountSpecificRule()

		if rule.AccountId != nil {
			continue
		}

		if rule.Forbids {
			return regources.SaleAccessDefinitionTypeWhitelist
		} else {
			return regources.SaleAccessDefinitionTypeBlacklist
		}
	}

	return regources.SaleAccessDefinitionTypeNone
}

func (c *saleHandler) convertSale(raw xdr.SaleEntry) (*history.Sale, error) {
	quoteAssets := make([]history.SaleQuoteAsset, 0, len(raw.QuoteAssets))
	for i := range raw.QuoteAssets {
		quoteAssets = append(quoteAssets, history.SaleQuoteAsset{
			Asset:          string(raw.QuoteAssets[i].QuoteAsset),
			Price:          regources.Amount(raw.QuoteAssets[i].Price),
			QuoteBalanceID: raw.QuoteAssets[i].QuoteBalance.AsString(),
			CurrentCap:     regources.Amount(raw.QuoteAssets[i].CurrentCap),
		})
	}

	saleType := raw.SaleTypeExt.SaleType

	return &history.Sale{
		ID:                uint64(raw.SaleId),
		OwnerAddress:      raw.OwnerId.Address(),
		BaseAsset:         string(raw.BaseAsset),
		DefaultQuoteAsset: string(raw.DefaultQuoteAsset),
		StartTime:         time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:           time.Unix(int64(raw.EndTime), 0).UTC(),
		SoftCap:           regources.Amount(raw.SoftCap),
		HardCap:           regources.Amount(raw.HardCap),
		Details:           internal.MarshalCustomDetails(raw.Details),
		QuoteAssets: history.SaleQuoteAssets{
			QuoteAssets: quoteAssets,
		},
		BaseCurrentCap: regources.Amount(raw.CurrentCapInBase),
		BaseHardCap:    regources.Amount(raw.MaxAmountToBeSold),
		SaleType:       saleType,
		// if sale still exists in core db - it is open
		State:   regources.SaleStateOpen,
		Version: int32(raw.Ext.V),
	}, nil
}
