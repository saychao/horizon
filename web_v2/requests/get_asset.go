package requests

import "net/http"

const (
	// IncludeTypeAssetOwner - defines if the asset owner should be included in the response
	IncludeTypeAssetOwner = "owner"
)

var includeTypeAssetAll = map[string]struct{}{
	IncludeTypeAssetOwner: {},
}

// GetAsset - represents params to be specified by user for Get Asset handler
type GetAsset struct {
	*base
	Code     string
	Includes struct {
		Owner bool `include:"owner"`
	}
}

// NewGetAsset returns new instance of GetAsset request
func NewGetAsset(r *http.Request) (*GetAsset, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAssetAll,
	})
	if err != nil {
		return nil, err
	}

	code := b.getString("code")

	return &GetAsset{
		base: b,
		Code: code,
	}, nil
}
