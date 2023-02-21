package operations

import (
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type createAtomicSwapAskRequestOpHandler struct {
	effectsProvider
}

// Details returns details about create atomic swap bid request operation
func (h *createAtomicSwapAskRequestOpHandler) Details(op rawOperation,
	opRes xdr.OperationResultTr) (history2.OperationDetails, error) {

	aSwapAskRequest := op.Body.MustCreateAtomicSwapAskRequestOp().Request

	quoteAssets := make([]regources.AssetPrice, 0, len(aSwapAskRequest.QuoteAssets))
	for _, quoteAsset := range aSwapAskRequest.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Asset: string(quoteAsset.QuoteAsset),
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	return history2.OperationDetails{
		Type: xdr.OperationTypeCreateAtomicSwapBidRequest,
		CreateAtomicSwapAskRequest: &history2.CreateAtomicSwapAskRequestDetails{
			Amount:      regources.Amount(aSwapAskRequest.Amount),
			BaseBalance: aSwapAskRequest.BaseBalance.AsString(),
			QuoteAssets: quoteAssets,
			Details:     internal.MarshalCustomDetails(aSwapAskRequest.CreatorDetails),
		},
	}, nil
}

// ParticipantsEffects returns source effect with `locked` effect
func (h *createAtomicSwapAskRequestOpHandler) ParticipantsEffects(opBody xdr.OperationBody,
	opRes xdr.OperationResultTr, sourceAccountID xdr.AccountId, _ []xdr.LedgerEntryChange,
) ([]history2.ParticipantEffect, error) {
	aSwapAskRequest := opBody.MustCreateAtomicSwapAskRequestOp().Request

	return []history2.ParticipantEffect{h.BalanceEffect(aSwapAskRequest.BaseBalance, &history2.Effect{
		Type: history2.EffectTypeLocked,
		Locked: &history2.BalanceChangeEffect{
			Amount: regources.Amount(aSwapAskRequest.Amount),
		},
	})}, nil
}
