package reviewablerequest

import (
	"strconv"

	"github.com/saychao/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateAtomicSwapRequest(histRequest history.AtomicSwap) (
	*regources.AtomicSwap, error,
) {
	return &regources.AtomicSwap{
		BidID:      strconv.FormatUint(histRequest.AskID, 10),
		BaseAmount: regources.Amount(histRequest.BaseAmount),
		QuoteAsset: histRequest.QuoteAsset,
	}, nil
}
