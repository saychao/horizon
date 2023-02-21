package exchange

import core "github.com/SafeRE-IT/horizon/db2/core2"

//go:generate mockery -case underscore -name assetProvider -inpkg -testonly
// assetProvider - helper interface to get assets from db
type assetProvider interface {
	//SelectByPolicy - selects slice of assets by policy mask
	SelectByPolicy(policy uint64) ([]core.Asset, error)
	// GetByCode - loads a row from `assets` found with matching code
	// returns nil, nil - if such asset doesn't exists
	GetByCode(code string) (*core.Asset, error)
	//SelectByAssets - selects slice of asset pairs where baseAsset in baseAssets and quoteAsset in quoteAssets
	SelectByAssets(baseAssets []string, quoteAssets []string) ([]core.AssetPair, error)
}
