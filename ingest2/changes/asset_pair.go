package changes

import (
	"time"

	history "github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

//go:generate mockery -case underscore -name assetPairStorage -inpkg -testonly
type assetPairStorage interface {
	InsertAssetPair(history.AssetPair) error
}

type assetPairHandler struct {
	storage assetPairStorage
}

func newAssetPairHandler(storage assetPairStorage) *assetPairHandler {
	return &assetPairHandler{
		storage: storage,
	}
}

func (p *assetPairHandler) Created(lc ledgerChange) error {
	assetPair := lc.LedgerChange.MustCreated().Data.MustAssetPair()
	if err := p.insert(assetPair, lc.LedgerCloseTime); err != nil {
		return errors.Wrap(err, "failed to insert from created")
	}
	return nil
}

func (p *assetPairHandler) Updated(lc ledgerChange) error {
	assetPair := lc.LedgerChange.MustUpdated().Data.MustAssetPair()
	if err := p.insert(assetPair, lc.LedgerCloseTime); err != nil {
		return errors.Wrap(err, "failed to insert from updated")
	}
	return nil
}

func (p *assetPairHandler) Removed(lc ledgerChange) error {
	// Nothing to do on asset pair remove
	return nil
}

func (p *assetPairHandler) insert(entry xdr.AssetPairEntry, ts time.Time) error {
	newAssetPair := history.NewAssetPair(string(entry.Base), string(entry.Quote), int64(entry.CurrentPrice), ts)
	err := p.storage.InsertAssetPair(newAssetPair)
	if err != nil {
		return errors.Wrap(err, "failed to insert asset pair")
	}
	return nil
}
