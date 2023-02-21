package changes

import (
	"encoding/json"
	"time"

	"github.com/SafeRE-IT/horizon/ingest2/storage"

	regources "gitlab.com/tokend/regources/generated"

	history "github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type pollStorage interface {
	//Inserts poll into DB
	Insert(poll history.Poll) error
	//Updates poll
	Update(poll history.Poll) error
	// SetState - sets state
	SetState(state regources.PollState) *storage.Poll
	// SetDetails - sets poll details
	SetDetails(details json.RawMessage) *storage.Poll
	// UpdateWhere - updates poll with values set by setters. Second parameter allows not to reset query after execution
	UpdateWhere(pollID uint64, shouldResetUpdater bool) error
}

type pollHandler struct {
	storage pollStorage
}

func newPollHandler(storage pollStorage) *pollHandler {
	return &pollHandler{
		storage: storage,
	}
}

// Created - handles creation of new poll
func (c *pollHandler) Created(lc ledgerChange) error {
	rawPoll := lc.LedgerChange.MustCreated().Data.MustPoll()

	poll, err := c.convertPoll(rawPoll)
	if err != nil {
		return errors.Wrap(err, "failed to convert poll", logan.F{
			"raw_poll":        rawPoll,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(poll)
	if err != nil {
		return errors.Wrap(err, "failed to insert poll into DB", logan.F{
			"poll": poll,
		})
	}

	return nil
}

// Removed - handles state of the poll due to it was removed
func (c *pollHandler) Removed(lc ledgerChange) error {
	pollID := uint64(lc.LedgerChange.MustRemoved().MustPoll().Id)
	managePollOp := lc.Operation.Body.MustManagePollOp()
	state, err := c.getPollState(managePollOp)
	if err != nil {
		return errors.Wrap(err, "failed to get poll state")
	}
	switch managePollOp.Data.Action {
	case xdr.ManagePollActionClose:
		{

			outcome := json.RawMessage(managePollOp.Data.ClosePollData.Details)
			err = c.storage.
				SetState(state).
				SetDetails(outcome).
				UpdateWhere(pollID, true)
		}
	case xdr.ManagePollActionCancel:
		{
			err = c.storage.
				SetState(state).
				UpdateWhere(pollID, true)
		}

	}

	if err != nil {
		return errors.Wrap(err, "failed to set poll outcome")
	}

	return nil
}

// Updated - handles update of the poll
func (c *pollHandler) Updated(lc ledgerChange) error {
	rawPoll := lc.LedgerChange.MustUpdated().Data.MustPoll()
	poll, err := c.convertPoll(rawPoll)
	if err != nil {
		return errors.Wrap(err, "failed to convert poll ", logan.F{
			"raw_poll":        rawPoll,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	managePollOp := lc.Operation.Body.MustManagePollOp()
	state, err := c.getPollState(managePollOp)
	if err != nil {
		return errors.Wrap(err, "failed to get poll state")
	}
	poll.State = state

	err = c.storage.Update(poll)
	if err != nil {
		return errors.Wrap(err, "failed to update poll", logan.F{
			"poll": poll,
		})
	}
	return nil
}

func (c *pollHandler) getPollState(op xdr.ManagePollOp) (regources.PollState, error) {
	var state regources.PollState
	switch op.Data.Action {
	case xdr.ManagePollActionCancel:
		state = regources.PollStateCancelled
	case xdr.ManagePollActionUpdateEndTime:
		state = regources.PollStateOpen
	case xdr.ManagePollActionClose:
		switch op.Data.MustClosePollData().Result {
		case xdr.PollResultFailed:
			state = regources.PollStateFailed
		case xdr.PollResultPassed:
			state = regources.PollStatePassed
		default:
			return state, errors.From(errors.New("Unexpected poll result"), logan.F{
				"poll_result": op.Data.MustClosePollData().Result,
			})
		}
	default:
		return state, errors.From(errors.New("Unexpected manage poll action"), logan.F{
			"action": op.Data.Action,
		})
	}
	return state, nil
}

func (c *pollHandler) convertPoll(raw xdr.PollEntry) (history.Poll, error) {

	return history.Poll{
		ID:               int64(raw.Id),
		OwnerID:          raw.OwnerId.Address(),
		ResultProviderID: raw.ResultProviderId.Address(),
		Data: history.PollData{
			Type: raw.Data.Type,
		},
		NumberOfChoices:          uint32(raw.NumberOfChoices),
		PermissionType:           uint32(raw.PermissionType),
		VoteConfirmationRequired: raw.VoteConfirmationRequired,
		StartTime:                time.Unix(int64(raw.StartTime), 0).UTC(),
		EndTime:                  time.Unix(int64(raw.EndTime), 0).UTC(),
		Details:                  internal.MarshalCustomDetails(raw.Details),
		State:                    regources.PollStateOpen,
	}, nil
}
