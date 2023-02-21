package horizon

import (
	"github.com/SafeRE-IT/horizon/render/hal"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/SafeRE-IT/horizon/resource"
	"gitlab.com/tokend/regources"
)

type AssetPairsAction struct {
	Action
	assetPairs []regources.AssetPair
}

func (action *AssetPairsAction) JSON() {
	action.Do(
		action.loadData,
		func() {
			hal.Render(action.W, action.assetPairs)
		},
	)
}

func (action *AssetPairsAction) loadData() {
	assets, err := action.CoreQ().AssetPairs().Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithStack(err).WithError(err).Error("Could not get asset from the database")
		return
	}

	action.assetPairs = make([]regources.AssetPair, len(assets))
	for i := range assets {
		action.assetPairs[i] = resource.PopulateAssetPair(assets[i])
	}
}
