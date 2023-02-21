// Package participants contains functions to derive a set of "participant"
// addresses for various data structures in the Stellar network's ledger.
package participants

import (
	"fmt"

	"github.com/saychao/horizon/db2/core"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// ForOperation returns all the participating accounts from the
// provided operation.

type Participant struct {
	AccountID xdr.AccountId
	BalanceID *xdr.BalanceId
	Details   interface{}
}

func ForOperation(
	DB *pgdb.DB,
	tx *xdr.Transaction,
	op *xdr.Operation,
	opResult xdr.OperationResultTr,
	ledgerChanges xdr.LedgerEntryChanges,
	ledger *core.LedgerHeader,
) (result []Participant, err error) {
	sourceParticipant := &Participant{}
	if op.SourceAccount != nil {
		sourceParticipant.AccountID = *op.SourceAccount
	} else {
		sourceParticipant.AccountID = tx.SourceAccount
	}
	switch op.Body.Type {
	case xdr.OperationTypeCreateAccount:
		result = append(result, Participant{op.Body.MustCreateAccountOp().Destination, nil, nil})
	case xdr.OperationTypeSetFees:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateWithdrawalRequest:
		createWithdrawalRequest := op.Body.MustCreateWithdrawalRequestOp()
		sourceParticipant.BalanceID = &createWithdrawalRequest.Request.Balance
	case xdr.OperationTypeManageBalance:
		manageBalanceOp := op.Body.MustManageBalanceOp()
		if sourceParticipant.AccountID.Address() != manageBalanceOp.Destination.Address() {
			result = append(result, Participant{manageBalanceOp.Destination, nil, nil})
		}
	case xdr.OperationTypeManageAsset:
		// the only direct participant is the source_accountWWW
	case xdr.OperationTypeManageLimits:
		// the only direct participant is the source_account, but I'm not sure
	case xdr.OperationTypeManageAssetPair:
		// the only direct participant is the source_account
	case xdr.OperationTypeManageOffer:
		manageOfferOp := op.Body.MustManageOfferOp()
		manageOfferResult := opResult.MustManageOfferResult()
		result = addMatchParticipants(result, sourceParticipant.AccountID, manageOfferOp.BaseBalance,
			manageOfferOp.QuoteBalance, manageOfferOp.IsBuy, manageOfferResult.Success)
		if manageOfferOp.IsBuy {
			sourceParticipant.BalanceID = &manageOfferOp.QuoteBalance
		} else {
			sourceParticipant.BalanceID = &manageOfferOp.BaseBalance
		}
		if len(result) != 0 {
			sourceParticipant = nil
		}
	case xdr.OperationTypeManageInvoiceRequest:
		manageInvoiceOp := op.Body.MustManageInvoiceRequestOp()
		switch manageInvoiceOp.Details.Action {
		case xdr.ManageInvoiceRequestActionCreate:
			result = append(result, Participant{manageInvoiceOp.Details.InvoiceRequest.Sender,
				&opResult.ManageInvoiceRequestResult.Success.Details.Response.SenderBalance, nil})
		case xdr.ManageInvoiceRequestActionRemove:
			sourceParticipant = nil
		}
	case xdr.OperationTypeManageContractRequest:
		manageContractOp := op.Body.MustManageContractRequestOp()
		switch manageContractOp.Details.Action {
		case xdr.ManageContractRequestActionCreate:
			result = append(result, Participant{manageContractOp.Details.MustCreateContractRequest().ContractRequest.Customer,
				nil, nil})
			result = append(result, Participant{manageContractOp.Details.MustCreateContractRequest().ContractRequest.Escrow,
				nil, nil})
		case xdr.ManageContractRequestActionRemove:
			sourceParticipant = nil
		}
	case xdr.OperationTypeManageContract:
		// the only direct participant is the source_account
	case xdr.OperationTypeReviewRequest:
		request := getReviewableRequestByID(uint64(op.Body.MustReviewRequestOp().RequestId), ledgerChanges)
		if request == nil || sourceParticipant.AccountID.Address() == request.Requestor.Address() {
			break
		}

		if request.Body.Type.ShortString() != xdr.ReviewableRequestTypeCreateAtomicSwapBid.ShortString() {
			result = append(result, Participant{
				AccountID: request.Requestor,
				BalanceID: nil,
				Details:   nil,
			})
			break
		}

		extendedResult, ok := opResult.MustReviewRequestResult().GetSuccess()
		if !ok {
			break
		}

		atomicSwapExtendedResult, ok := extendedResult.TypeExt.GetAtomicSwapBidExtended()
		if !ok {
			break
		}

		result = append(result, Participant{
			AccountID: atomicSwapExtendedResult.BidOwnerId,
			BalanceID: &atomicSwapExtendedResult.BidOwnerBaseBalanceId,
			Details:   nil,
		})

		result = append(result, Participant{
			AccountID: atomicSwapExtendedResult.AskOwnerId,
			BalanceID: &atomicSwapExtendedResult.AskOwnerBaseBalanceId,
			Details:   nil,
		})
	case xdr.OperationTypeCreatePreissuanceRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateIssuanceRequest:
		manageIssuanceRequest := op.Body.MustCreateIssuanceRequestOp()
		manageIssuanceResult := opResult.MustCreateIssuanceRequestResult()
		result = append(result, Participant{manageIssuanceResult.MustSuccess().Receiver,
			&manageIssuanceRequest.Request.Receiver, nil})
	case xdr.OperationTypeCreateSaleRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCheckSaleState:
		manageOfferResult := opResult.MustCheckSaleStateResult()
		saleClosed := manageOfferResult.Success.Effect.SaleClosed
		if saleClosed == nil {
			break
		}

		for i := range saleClosed.Results {
			result = addMatchParticipants(result, saleClosed.SaleOwner, saleClosed.Results[i].SaleBaseBalance,
				saleClosed.Results[i].SaleQuoteBalance, false, &saleClosed.Results[i].SaleDetails)
		}

		sourceParticipant = nil
	case xdr.OperationTypeManageExternalSystemAccountIdPoolEntry:
		// the only direct participant is the source_account
	case xdr.OperationTypeBindExternalSystemAccountId:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateAmlAlert:
		// TODO add participant
	case xdr.OperationTypePayment:
		paymentOpV2 := op.Body.MustPaymentOp()
		paymentV2Response := opResult.MustPaymentResult().MustPaymentResponse()

		result = append(result, Participant{paymentV2Response.Destination, &paymentV2Response.DestinationBalanceId, nil})
		sourceParticipant.BalanceID = &paymentOpV2.SourceBalanceId
	case xdr.OperationTypeManageSale:
		// the only direct participant is the source_account
	case xdr.OperationTypeManageKeyValue:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateManageLimitsRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateAtomicSwapAskRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCancelAtomicSwapAsk:
		// the only direct participant is the source_account
	case xdr.OperationTypeStamp:
		// the only direct participant is the source_account
	case xdr.OperationTypeLicense:
		// the only direct participant is the source_account
	case xdr.OperationTypeManagePoll:
	case xdr.OperationTypeManageCreatePollRequest:
	case xdr.OperationTypeManageVote:
	case xdr.OperationTypeRemoveAsset:
		// the only direct participant is the source_account
	case xdr.OperationTypeCreateAtomicSwapBidRequest:
		// FIXME !!!!!
		//tx.SourceAccount
		if sourceParticipant.AccountID.Address() != opResult.MustCreateAtomicSwapBidRequestResult().Success.AskOwnerId.Address() {
			result = append(result, Participant{opResult.MustCreateAtomicSwapBidRequestResult().Success.AskOwnerId, nil, nil})
		}
	case xdr.OperationTypeCancelSaleRequest:
		// the only direct participant is the source_account
	case xdr.OperationTypeCancelChangeRoleRequest:
	case xdr.OperationTypeManageAccountRule:
	case xdr.OperationTypeManageSignerRule:
	case xdr.OperationTypeManageSigner:
	case xdr.OperationTypeManageAccountRole:
	case xdr.OperationTypeManageSignerRole:
	case xdr.OperationTypeCreateChangeRoleRequest:
		updateKYCRequestData := op.Body.MustCreateChangeRoleRequestOp()
		if sourceParticipant.AccountID.Address() != updateKYCRequestData.DestinationAccount.Address() {
			result = append(result, Participant{
				AccountID: updateKYCRequestData.DestinationAccount,
				BalanceID: nil,
			})
		}
	case xdr.OperationTypeManageAccountSpecificRule:
	case xdr.OperationTypeRemoveAssetPair:
	case xdr.OperationTypeCreateKycRecoveryRequest:
	case xdr.OperationTypeInitiateKycRecovery:
	case xdr.OperationTypeCreatePaymentRequest:
	case xdr.OperationTypeCreateManageOfferRequest:
	case xdr.OperationTypeOpenSwap:
	case xdr.OperationTypeCloseSwap:
	case xdr.OperationTypeCreateData:
	case xdr.OperationTypeRemoveData:
	case xdr.OperationTypeUpdateData:
	case xdr.OperationTypeUpdateDataOwner:
	case xdr.OperationTypeCreateRedemptionRequest:
		// TODO add participant
	case xdr.OperationTypeCreateDataCreationRequest:
	case xdr.OperationTypeCreateDataUpdateRequest:
	case xdr.OperationTypeCreateDataOwnerUpdateRequest:
	case xdr.OperationTypeCreateDataRemoveRequest:
	case xdr.OperationTypeCancelDataCreationRequest:
	case xdr.OperationTypeCancelDataUpdateRequest:
	case xdr.OperationTypeCancelDataOwnerUpdateRequest:
	case xdr.OperationTypeCancelDataRemoveRequest:
	case xdr.OperationTypeCreateDeferredPaymentCreationRequest:
	case xdr.OperationTypeCancelDeferredPaymentCreationRequest:
	case xdr.OperationTypeCreateCloseDeferredPaymentRequest:
	case xdr.OperationTypeCancelCloseDeferredPaymentRequest:
	default:
		err = fmt.Errorf("unknown operation type: %s", op.Body.Type)
	}

	if sourceParticipant != nil {
		result = append(result, *sourceParticipant)
	}
	return
}

func getReviewableRequestByID(id uint64, ledgerChanges xdr.LedgerEntryChanges) *xdr.ReviewableRequestEntry {
	for i := range ledgerChanges {
		ledgerChange := ledgerChanges[i]
		var reviewableRequest *xdr.ReviewableRequestEntry
		switch ledgerChange.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			reviewableRequest = ledgerChange.Created.Data.ReviewableRequest
		case xdr.LedgerEntryChangeTypeUpdated:
			reviewableRequest = ledgerChange.Updated.Data.ReviewableRequest
		default:
			continue
		}

		if reviewableRequest == nil {
			continue
		}

		if uint64(reviewableRequest.RequestId) == id {
			return reviewableRequest
		}
	}

	return nil
}

func addMatchParticipants(participants []Participant, offerSourceID xdr.AccountId, baseBalanceID xdr.BalanceId,
	quoteBalanceID xdr.BalanceId, isBuy bool, result *xdr.ManageOfferSuccessResult) []Participant {
	if result == nil || len(result.OffersClaimed) == 0 {
		return participants
	}

	matchesByBalance := NewMatchesDetailsByBalance()

	for _, offerClaimed := range result.OffersClaimed {

		claimedOfferMatch := NewMatch(offerClaimed.BaseAmount, offerClaimed.QuoteAmount, offerClaimed.BFeePaid, offerClaimed.CurrentPrice)
		matchesByBalance.Add(offerClaimed.BAccountId, []xdr.BalanceId{offerClaimed.BaseBalance, offerClaimed.QuoteBalance},
			result.BaseAsset, result.QuoteAsset, !isBuy, claimedOfferMatch)

		offerMatch := NewMatch(offerClaimed.BaseAmount, offerClaimed.QuoteAmount, offerClaimed.AFeePaid, offerClaimed.CurrentPrice)
		matchesByBalance.Add(offerSourceID, []xdr.BalanceId{baseBalanceID, quoteBalanceID}, result.BaseAsset, result.QuoteAsset, isBuy, offerMatch)

	}

	return matchesByBalance.ToParticipants(participants)
}

// ForTransaction returns all the participating accounts from the provided
// transaction.
func ForTransaction(
	DB *pgdb.DB,
	tx *xdr.Transaction,
	opResults []xdr.OperationResult,
	meta *xdr.TransactionMeta,
	ledgerChanges *xdr.LedgerEntryChanges,
	ledger *core.LedgerHeader,
) (result []xdr.AccountId, err error) {

	result = append(result, tx.SourceAccount)

	p, err := forMeta(meta)
	if err != nil {
		return
	}
	result = append(result, p...)

	p, err = forChanges(ledgerChanges)
	if err != nil {
		return
	}
	result = append(result, p...)

	for i := range tx.Operations {
		participants, err := ForOperation(DB, tx, &tx.Operations[i], *opResults[i].Tr, nil, ledger)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get participants for operation")
		}
		for _, participant := range participants {
			result = append(result, participant.AccountID)
		}
	}

	result = dedupe(result)
	return
}

// dedupe remove any duplicate ids from `in`
func dedupe(in []xdr.AccountId) (out []xdr.AccountId) {
	set := map[string]xdr.AccountId{}
	for _, id := range in {
		set[id.Address()] = id
	}

	for _, id := range set {
		out = append(out, id)
	}
	return
}

func forChanges(
	changes *xdr.LedgerEntryChanges,
) (result []xdr.AccountId, err error) {

	for _, c := range *changes {
		var account *xdr.AccountId

		switch c.Type {
		case xdr.LedgerEntryChangeTypeCreated:
			account = forLedgerEntry(c.MustCreated())
		case xdr.LedgerEntryChangeTypeRemoved:
			account = forLedgerKey(c.MustRemoved())
		case xdr.LedgerEntryChangeTypeUpdated:
			account = forLedgerEntry(c.MustUpdated())
		case xdr.LedgerEntryChangeTypeState:
			account = forLedgerEntry(c.MustState())
		default:
			err = fmt.Errorf("Unknown change type: %s", c.Type)
			return
		}

		if account != nil {
			result = append(result, *account)
		}
	}

	return
}

func forLedgerEntry(le xdr.LedgerEntry) *xdr.AccountId {
	if le.Data.Type != xdr.LedgerEntryTypeAccount {
		return nil
	}
	aid := le.Data.MustAccount().AccountId
	return &aid
}

func forLedgerKey(lk xdr.LedgerKey) *xdr.AccountId {
	if lk.Type != xdr.LedgerEntryTypeAccount {
		return nil
	}
	aid := lk.MustAccount().AccountId
	return &aid
}

func forMeta(
	meta *xdr.TransactionMeta,
) (result []xdr.AccountId, err error) {

	if meta.Operations == nil {
		return
	}

	for _, op := range *meta.Operations {
		var acc []xdr.AccountId
		acc, err = forChanges(&op.Changes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get ledger changes")
		}

		result = append(result, acc...)
	}

	return
}
