package resources

import (
	"fmt"
	"time"

	"github.com/saychao/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

// NewAccount - creates new instance of account
func NewAtomicSwapAsk(core core2.AtomicSwapAsk) regources.AtomicSwapAsk {
	return regources.AtomicSwapAsk{
		Key: NewAtomicSwapAskKey(uint64(core.AskID)),
		Attributes: regources.AtomicSwapAskAttributes{
			AvailableAmount: regources.Amount(core.AvailableAmount),
			LockedAmount:    regources.Amount(core.LockedAmount),
			CreatedAt:       time.Unix(int64(core.CreatedAt), 0).UTC(),
			IsCanceled:      core.IsCanceled,
			Details:         core.Details,
		},
	}
}

// NewAccountKey - creates account key from address
func NewAtomicSwapAskKey(id uint64) regources.Key {
	return regources.Key{
		ID:   fmt.Sprint(id),
		Type: regources.ATOMIC_SWAP_ASK,
	}
}

func NewAtomicSwapAskQuoteAssetKey(quoteAsset string, askID uint64) regources.Key {
	return regources.Key{
		// TODO: Use artificial ID
		ID:   fmt.Sprintf("%s:%d", quoteAsset, askID),
		Type: regources.ATOMIC_SWAP_QUOTE_ASSETS,
	}
}

func NewAtomicSwapAskQuoteAsset(raw core2.AtomicSwapQuoteAsset) regources.AtomicSwapQuoteAsset {
	return regources.AtomicSwapQuoteAsset{
		Key: NewAtomicSwapAskQuoteAssetKey(raw.QuoteAsset, uint64(raw.AskID)),
		Attributes: regources.AtomicSwapQuoteAssetAttributes{
			Price:      regources.Amount(raw.Price),
			QuoteAsset: raw.QuoteAsset,
		},
	}
}
