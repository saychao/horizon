package reviewablerequest

import (
	"github.com/SafeRE-IT/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateAtomicSwapAskCreationRequest(histRequest history.AtomicSwapAskCreation) (
	*regources.AtomicSwapBidCreation, error,
) {
	return &regources.AtomicSwapBidCreation{
		BaseBalance: histRequest.BaseBalance,
		BaseAmount:  regources.Amount(histRequest.BaseAmount),
		QuoteAssets: histRequest.QuoteAssets,
		Details:     histRequest.Details,
	}, nil
}
