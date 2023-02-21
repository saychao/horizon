package changes

import (
	"encoding/hex"

	"github.com/saychao/horizon/ingest2/internal"
	regources "gitlab.com/tokend/regources/generated"

	history "github.com/saychao/horizon/db2/history2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type swapStorage interface {
	//Inserts swap into DB
	Insert(swap history.Swap) error
	// SetState - sets state
	SetState(id int64, state regources.SwapState) error

	SetSecret(id int64, secret string) error
}

type swapHandler struct {
	storage swapStorage
}

func newSwapHandler(storage swapStorage) *swapHandler {
	return &swapHandler{
		storage: storage,
	}
}

// Created - handles creation of new swap
func (c *swapHandler) Created(lc ledgerChange) error {
	rawSwap := lc.LedgerChange.MustCreated().Data.MustSwap()
	op := lc.Operation.Body.MustOpenSwapOp()
	opRes := lc.OperationResult.MustOpenSwapResult()
	swap, err := c.convertSwap(rawSwap, op, opRes)
	if err != nil {
		return errors.Wrap(err, "failed to convert swap", logan.F{
			"swap_id":         rawSwap.Id,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(swap)
	if err != nil {
		return errors.Wrap(err, "failed to insert swap into DB", logan.F{
			"swap": swap,
		})
	}

	return nil
}

// Removed - handles state of the swap due to it was removed
func (c *swapHandler) Removed(lc ledgerChange) error {
	swapID := int64(lc.LedgerChange.MustRemoved().MustSwap().Id)
	closeSwapOp := lc.Operation.Body.MustCloseSwapOp()
	res := lc.OperationResult.MustCloseSwapResult()
	state, err := c.getSwapState(closeSwapOp, res)
	if err != nil {
		return errors.Wrap(err, "failed to get swap state")
	}

	err = c.storage.SetState(swapID, state)
	if err != nil {
		return errors.Wrap(err, "failed to set swap state")
	}

	if state == regources.SwapStateClosed {
		secret := hex.EncodeToString(closeSwapOp.Secret[:])
		err = c.storage.SetSecret(swapID, secret)
		if err != nil {
			return errors.Wrap(err, "failed to set swap secret")
		}
	}

	return nil
}

func (c *swapHandler) getSwapState(op xdr.CloseSwapOp, res xdr.CloseSwapResult) (regources.SwapState, error) {
	var state regources.SwapState
	switch res.MustSuccess().Effect {
	case xdr.CloseSwapEffectClosed:
		state = regources.SwapStateClosed
	case xdr.CloseSwapEffectCancelled:
		state = regources.SwapStateCanceled
	default:
		return state, errors.From(errors.New("Unexpected close swap effect"), logan.F{
			"effect": res.MustSuccess().Effect,
		})
	}
	return state, nil
}

func (c *swapHandler) convertSwap(raw xdr.SwapEntry, op xdr.OpenSwapOp, res xdr.OpenSwapResult) (history.Swap, error) {
	success := res.MustSuccess()

	secretHash := hex.EncodeToString(op.SecretHash[:])

	return history.Swap{
		ID:                    int64(raw.Id),
		SourceAccount:         raw.Source.Address(),
		CreatedAt:             internal.TimeFromXdr(xdr.Uint64(raw.CreatedAt)),
		LockTime:              internal.TimeFromXdr(xdr.Uint64(raw.LockTime)),
		DestinationAccount:    success.Destination.Address(),
		SourceBalance:         raw.SourceBalance.AsString(),
		DestinationBalance:    success.DestinationBalance.AsString(),
		SecretHash:            secretHash,
		Amount:                uint64(op.Amount),
		Asset:                 string(success.Asset),
		SourceFixedFee:        uint64(success.ActualSourceFee.Fixed),
		SourcePercentFee:      uint64(success.ActualSourceFee.Percent),
		DestinationFixedFee:   uint64(success.ActualDestinationFee.Fixed),
		DestinationPercentFee: uint64(success.ActualDestinationFee.Percent),
		Details:               internal.MarshalCustomDetails(raw.Details),
		State:                 regources.SwapStateOpen,
	}, nil
}
