package changes

import (
	history "github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

//go:generate mockery -case underscore -name assetPairStorage -inpkg -testonly
type assetStorage interface {
	Insert(history.Asset) error
	Update(asset history.Asset) error
	SetState(code string, state regources.AssetState) error
}

type assetHandler struct {
	storage assetStorage
}

func newAssetHandler(storage assetStorage) *assetHandler {
	return &assetHandler{
		storage: storage,
	}
}

func (h *assetHandler) Created(lc ledgerChange) error {
	rawAsset := lc.LedgerChange.MustCreated().Data.MustAsset()
	asset := h.convertAsset(rawAsset)
	if err := h.storage.Insert(asset); err != nil {
		return errors.Wrap(err, "failed to insert from created")
	}
	return nil
}

func (h *assetHandler) Updated(lc ledgerChange) error {
	rawAsset := lc.LedgerChange.MustUpdated().Data.MustAsset()
	asset := h.convertAsset(rawAsset)
	if err := h.storage.Update(asset); err != nil {
		return errors.Wrap(err, "failed to update from updated")
	}
	return nil
}

func (h *assetHandler) Stated(lc ledgerChange) error {
	op := lc.Operation.Body
	if op.Type == xdr.OperationTypeRemoveAsset {
		if err := h.storage.SetState(string(op.RemoveAssetOp.Code), regources.AssetStateDeleted); err != nil {
			return errors.Wrap(err, "failed to set state from stated")
		}
	}

	return nil
}

func (h *assetHandler) convertAsset(raw xdr.AssetEntry) history.Asset {
	return history.Asset{
		Code:                   string(raw.Code),
		Owner:                  raw.Owner.Address(),
		PreIssuanceAssetSigner: raw.PreissuedAssetSigner.Address(),
		Details:                []byte(raw.Details),
		MaxIssuanceAmount:      uint64(raw.MaxIssuanceAmount),
		AvailableForIssuance:   uint64(raw.AvailableForIssueance),
		Issued:                 uint64(raw.Issued),
		PendingIssuance:        uint64(raw.PendingIssuance),
		Policies:               uint32(raw.Policies),
		TrailingDigits:         uint32(raw.TrailingDigitsCount),
		Type:                   uint64(raw.Type),
		State:                  regources.AssetStateActive,
	}
}
