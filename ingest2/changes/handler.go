package changes

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	core "github.com/SafeRE-IT/horizon/db2/core2"
)

type creatable interface {
	Created(change ledgerChange) error
}
type updatable interface {
	Updated(change ledgerChange) error
}
type removable interface {
	Removed(change ledgerChange) error
}
type statable interface {
	Stated(change ledgerChange) error
}

// Handler - handles ledger changes to populate changes for entries
type Handler struct {
	Create map[xdr.LedgerEntryType]creatable
	Update map[xdr.LedgerEntryType]updatable
	Remove map[xdr.LedgerEntryType]removable
	State  map[xdr.LedgerEntryType]statable
}

//NewHandler - returns new instance of handler
func NewHandler(account accountStorage,
	balance balanceStorage,
	request reviewableRequestStorage,
	sale saleStorage,
	assetPair assetPairStorage,
	poll pollStorage,
	vote voteStorage,
	accountSpecificRule accountSpecificRuleStorage,
	accountStatus accountStatusStorage,
	swap swapStorage,
	asset assetStorage,
	data dataStorage,
	def deferredPaymentStorage,
) *Handler {

	reviewRequestHandlerInst := newReviewableRequestHandler(request, balance, accountStatus)
	saleHandlerInst := newSaleHandler(sale, accountSpecificRule)
	assetPairHandler := newAssetPairHandler(assetPair)
	pollHandlerInst := newPollHandler(poll)
	voteHandlerInst := newVoteHandler(vote)
	accountSpecificRuleHandlerInst := newAccountSpecificRuleHandler(accountSpecificRule)
	signerHandlerInst := newSignerHandler(accountStatus)
	swapHandlerInst := newSwapHandler(swap)
	assetHandlerInst := newAssetHandler(asset)
	dataHandlerInst := newDataHandler(data)
	deferredPaymentHandlerInst := newDeferredPaymentHandler(def)
	return &Handler{
		Create: map[xdr.LedgerEntryType]creatable{
			xdr.LedgerEntryTypeAccount:             newAccountHandler(account),
			xdr.LedgerEntryTypeBalance:             newBalanceHandler(account, balance),
			xdr.LedgerEntryTypeReviewableRequest:   reviewRequestHandlerInst,
			xdr.LedgerEntryTypeSale:                saleHandlerInst,
			xdr.LedgerEntryTypeAssetPair:           assetPairHandler,
			xdr.LedgerEntryTypePoll:                pollHandlerInst,
			xdr.LedgerEntryTypeVote:                voteHandlerInst,
			xdr.LedgerEntryTypeAccountSpecificRule: accountSpecificRuleHandlerInst,
			xdr.LedgerEntryTypeSigner:              signerHandlerInst,
			xdr.LedgerEntryTypeSwap:                swapHandlerInst,
			xdr.LedgerEntryTypeAsset:               assetHandlerInst,
			xdr.LedgerEntryTypeData:                dataHandlerInst,
			xdr.LedgerEntryTypeDeferredPayment:     deferredPaymentHandlerInst,
		},
		Update: map[xdr.LedgerEntryType]updatable{
			xdr.LedgerEntryTypeReviewableRequest: reviewRequestHandlerInst,
			xdr.LedgerEntryTypeSale:              saleHandlerInst,
			xdr.LedgerEntryTypeAssetPair:         assetPairHandler,
			xdr.LedgerEntryTypePoll:              pollHandlerInst,
			xdr.LedgerEntryTypeAsset:             assetHandlerInst,
			xdr.LedgerEntryTypeData:              dataHandlerInst,
			xdr.LedgerEntryTypeDeferredPayment:   deferredPaymentHandlerInst,
		},
		Remove: map[xdr.LedgerEntryType]removable{
			xdr.LedgerEntryTypeReviewableRequest:   reviewRequestHandlerInst,
			xdr.LedgerEntryTypeSale:                saleHandlerInst,
			xdr.LedgerEntryTypePoll:                pollHandlerInst,
			xdr.LedgerEntryTypeVote:                voteHandlerInst,
			xdr.LedgerEntryTypeAssetPair:           assetPairHandler,
			xdr.LedgerEntryTypeAccountSpecificRule: accountSpecificRuleHandlerInst,
			xdr.LedgerEntryTypeSwap:                swapHandlerInst,
			xdr.LedgerEntryTypeData:                dataHandlerInst,
			xdr.LedgerEntryTypeDeferredPayment:     deferredPaymentHandlerInst,
		},
		State: map[xdr.LedgerEntryType]statable{
			xdr.LedgerEntryTypeReviewableRequest: reviewRequestHandlerInst,
			xdr.LedgerEntryTypeAsset:             assetHandlerInst,
		},
	}
}

// Handle - processes all the ledger changes for specified ledger
func (h *Handler) Handle(header *core.LedgerHeader, txs []core.Transaction) error {
	for txI := range txs {
		tx := txs[txI]
		ops := tx.ResultMeta.MustOperations()
		for opI := range ops {
			op := ops[opI]
			for changeI := range op.Changes {
				change := op.Changes[changeI]
				err := h.handle(ledgerChange{
					LedgerSeq:       header.Sequence,
					LedgerCloseTime: time.Unix(header.CloseTime, 0).UTC(),
					LedgerChange:    change,
					Operation:       &tx.Envelope.Tx.Operations[opI],
					OperationResult: tx.Result.Result.Result.MustResults()[opI].Tr,
					OperationIndex:  uint32(opI),
					OpChanges:       op.Changes,
					TxIndex:         uint32(tx.Index),
				})

				if err != nil {
					return errors.Wrap(err, "failed to process ledger change", logan.F{
						"ledger_seq": header.Sequence,
						"tx_i":       txI,
						"op_i":       opI,
						"change_i":   changeI,
						"change":     change,
					})
				}
			}
		}
	}

	return nil
}

// handle - tries to find corresponding ledger change handler and handle it.
func (h *Handler) handle(lc ledgerChange) error {
	switch lc.LedgerChange.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		return h.created(lc)
	case xdr.LedgerEntryChangeTypeUpdated:
		return h.updated(lc)
	case xdr.LedgerEntryChangeTypeRemoved:
		return h.removed(lc)
	case xdr.LedgerEntryChangeTypeState:
		return h.stated(lc)
	default:
		return errors.From(errors.New("Unrecognized ledger entry change type"), logan.F{
			"change_type": lc.LedgerChange.Type,
		})
	}
}

func (h *Handler) created(lc ledgerChange) error {
	handler, ok := h.Create[lc.LedgerChange.Created.Data.Type]
	if !ok {
		return nil
	}

	return handler.Created(lc)
}

func (h *Handler) updated(lc ledgerChange) error {
	handler, ok := h.Update[lc.LedgerChange.Updated.Data.Type]
	if !ok {
		return nil
	}

	return handler.Updated(lc)
}

func (h *Handler) removed(lc ledgerChange) error {
	handler, ok := h.Remove[lc.LedgerChange.Removed.Type]
	if !ok {
		return nil
	}

	return handler.Removed(lc)
}

func (h *Handler) stated(lc ledgerChange) error {
	handler, ok := h.State[lc.LedgerChange.State.Data.Type]
	if !ok {
		return nil
	}

	return handler.Stated(lc)
}

//Name - name of the handler
func (h *Handler) Name() string {
	return "ledger_changes_handler"
}
