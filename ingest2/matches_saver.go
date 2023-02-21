package ingest2

import (
	core "github.com/saychao/horizon/db2/core2"
	"github.com/saychao/horizon/db2/history2"
	"github.com/saychao/horizon/ingest2/generator"
	"github.com/saychao/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type matchesStorage interface {
	Insert(matches []history2.Match) error
}

// MatchesSaver - converts claimed offer atoms into horizon matches and stores them to db
type MatchesSaver struct {
	storage matchesStorage
}

// NewMatchesSaver - creates new instance of MatchesSaver
func NewMatchesSaver(storage matchesStorage) *MatchesSaver {
	return &MatchesSaver{
		storage: storage,
	}
}

// Handle - converts claimed offer atoms into horizon matches and stores them to db
func (h *MatchesSaver) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	matchIDGen := generator.NewIDI32(header.Sequence)
	opIDGen := generator.NewIDI32(header.Sequence)

	var ledgerMatches []history2.Match

	for _, tx := range txs {
		ops := tx.Envelope.Tx.Operations
		for opI, op := range ops {
			if op.Body.Type != xdr.OperationTypeManageOffer {
				continue
			}

			opID := opIDGen.Next()
			opResult := tx.Result.Result.Result.MustResults()[opI].MustTr()
			opMatches := getOpMatches(opResult, opID, header.CloseTime, matchIDGen)

			ledgerMatches = append(ledgerMatches, opMatches...)
		}
	}

	if len(ledgerMatches) == 0 {
		return nil
	}

	err := h.storage.Insert(ledgerMatches)
	if err != nil {
		return errors.Wrap(err, "failed to insert matches")
	}

	return nil
}

func getOpMatches(opResult xdr.OperationResultTr, opID int64, ledgerCloseTime int64, matchIDGen *generator.ID) []history2.Match {
	manageOfferResult := opResult.MustManageOfferResult().MustSuccess()

	var result []history2.Match

	for _, atom := range manageOfferResult.OffersClaimed {
		var ok bool
		if result, ok = trySquash(result, atom); !ok {
			result = append(result, history2.NewMatch(
				matchIDGen.Next(),
				opID,
				manageOfferResult.BaseAsset,
				manageOfferResult.QuoteAsset,
				internal.TimeFromXdr(xdr.Uint64(ledgerCloseTime)),
				atom,
			))
		}
	}

	return result
}

func trySquash(matches []history2.Match, atom xdr.ClaimOfferAtom) (m []history2.Match, ok bool) {
	for i, match := range matches {
		if match.Price == uint64(atom.CurrentPrice) {
			match.BaseAmount += uint64(atom.BaseAmount)
			match.QuoteAmount += uint64(atom.QuoteAmount)

			matches[i] = match

			return matches, true
		}
	}

	return matches, false
}

// Name - returns name of the handler
func (h *MatchesSaver) Name() string {
	return "matches_saver"
}
