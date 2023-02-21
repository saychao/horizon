package changes

import (
	history "github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/generator"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type voteStorage interface {
	//Inserts vote into DB
	Insert(vote history.Vote) error
	//Updates vote
	Update(vote history.Vote) error

	Remove(voterID string, pollID uint64) error
}

type voteHandler struct {
	storage voteStorage
}

func newVoteHandler(storage voteStorage) *voteHandler {
	return &voteHandler{
		storage: storage,
	}
}

// Created - handles creation of new vote
func (c *voteHandler) Created(lc ledgerChange) error {
	rawVote := lc.LedgerChange.MustCreated().Data.MustVote()
	//TX count is restricted to 500, so we must be safe
	// Op count in tx restricted to 100
	txInd := uint16(lc.TxIndex)
	opInd := uint16(lc.OperationIndex)
	voteID := generator.MakeIDUint16(lc.LedgerSeq, txInd, opInd)
	vote, err := c.convertVote(rawVote, voteID)
	if err != nil {
		return errors.Wrap(err, "failed to convert vote", logan.F{
			"raw_vote":        rawVote,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(*vote)
	if err != nil {
		return errors.Wrap(err, "failed to insert vote into DB", logan.F{
			"vote": vote,
		})
	}

	return nil
}

// Removed - handles state of the vote due to it was removed
func (c *voteHandler) Removed(lc ledgerChange) error {
	switch lc.Operation.Body.Type {
	case xdr.OperationTypeManageVote:
		pollID := uint64(lc.LedgerChange.MustRemoved().MustVote().PollId)
		voterID := lc.LedgerChange.MustRemoved().MustVote().VoterId
		err := c.storage.Remove(voterID.Address(), pollID)
		if err != nil {
			return errors.Wrap(err, "failed to remove vote")
		}
	case xdr.OperationTypeManagePoll:
		managePollOp := lc.Operation.Body.MustManagePollOp()
		switch managePollOp.Data.Action {
		case xdr.ManagePollActionUpdateEndTime:
		case xdr.ManagePollActionClose:
		case xdr.ManagePollActionCancel:
		default:
			return errors.From(
				errors.New("Unexpected manage poll action"),
				logan.F{
					"action": managePollOp.Data.Action,
				})
		}
	default:
		return errors.From(
			errors.New("Unexpected operation triggered VoteEntry remove"),
			logan.F{
				"operation_type": lc.Operation.Body.Type,
			})
	}
	return nil
}

// Updated - handles update of the vote
func (c *voteHandler) Updated(lc ledgerChange) error {
	rawVote := lc.LedgerChange.MustUpdated().Data.MustVote()
	vote, err := c.convertVote(rawVote, 0)
	if err != nil {
		return errors.Wrap(err, "failed to convert vote ", logan.F{
			"raw_vote":        rawVote,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Update(*vote)
	if err != nil {
		return errors.Wrap(err, "failed to update vote", logan.F{
			"vote": vote,
		})
	}
	return nil
}

func (c *voteHandler) convertVote(raw xdr.VoteEntry, ID int64) (*history.Vote, error) {
	choice := uint64(raw.Data.MustSingle().Choice)
	return &history.Vote{
		ID:      ID,
		VoterID: raw.VoterId.Address(),
		PollID:  int64(raw.PollId),
		VoteData: history.VoteData{
			PollType:     raw.Data.PollType,
			SingleChoice: &choice,
		},
	}, nil
}
