package requests

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	// IncludeTypeAssetPairBaseAsset - defines if the base asset should be included in the response
	IncludeTypeAssetPairBaseAsset = "base_asset"
	// IncludeTypeAssetPairQuoteAsset - defines if the quote asset should be included in the response
	IncludeTypeAssetPairQuoteAsset = "quote_asset"
)

var includeTypeAssetPairAll = map[string]struct{}{
	IncludeTypeAssetPairBaseAsset:  {},
	IncludeTypeAssetPairQuoteAsset: {},
}

// GetAssetPair - represents params to be specified by user for getAssetPair handler
type GetAssetPair struct {
	*base

	BaseAsset  string
	QuoteAsset string
	Includes   struct {
		BaseAsset  bool `include:"base_asset"`
		QuoteAsset bool `include:"quote_asset"`
	}
}

// NewGetAssetPair returns new instance of GetAssetPair request
func NewGetAssetPair(r *http.Request) (*GetAssetPair, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAssetPairAll,
	})
	if err != nil {
		return nil, err
	}

	request := GetAssetPair{
		base: b,
	}

	err = request.getAssetCodes("id")
	if err != nil {
		return nil, err
	}

	return &request, nil
}

// getAssetCodes receives the compound ID in `base:quote` format and
// returns separately both codes or validation error if id is invalid
func (r *GetAssetPair) getAssetCodes(param string) error {
	id := r.getString(param)

	codes := strings.Split(id, ":")
	if len(codes) != 2 || codes[0] == "" || codes[1] == "" {
		return validation.Errors{
			"id": errors.New("should be in `base:quote` format"),
		}
	}

	if codes[0] == codes[1] {
		return validation.Errors{
			"id": errors.New("can't contain equal base:quote values"),
		}
	}

	r.BaseAsset = codes[0]
	r.QuoteAsset = codes[1]

	return nil
}
