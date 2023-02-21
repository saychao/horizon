package changes

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/SafeRE-IT/horizon/db2"

	history "github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"github.com/SafeRE-IT/horizon/utf8"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

var errUnknownRemoveReason = errors.New("request was removed due to unknown reason")
var removeOnKYCRecoveryInit = "New KYC recovery was initiated"

type reviewableRequestStorage interface {
	//Inserts Reviewable request into DB
	Insert(request history.ReviewableRequest) error
	//Updates Reviewable request
	Update(request history.ReviewableRequest) error
	//Approves reviewable request
	Approve(id uint64) error
	//PermanentReject - rejects permanently reviewable request
	PermanentReject(id uint64, rejectReason string) error
	//Cancel - cancels reviewable request
	Cancel(id uint64) error
}

type balanceProvider interface {
	// MustBalance returns history balance struct for specific balance id
	MustBalance(balanceID xdr.BalanceId) history.Balance
}

type reviewableRequestHandler struct {
	storage  reviewableRequestStorage
	balances balanceProvider
	accounts accountStatusStorage
}

func newReviewableRequestHandler(storage reviewableRequestStorage, balances balanceProvider, accounts accountStatusStorage) *reviewableRequestHandler {
	return &reviewableRequestHandler{
		storage:  storage,
		balances: balances,
		accounts: accounts,
	}
}

// Created - handles creation of new reviewable request
func (c *reviewableRequestHandler) Created(lc ledgerChange) error {
	reviewableRequest := lc.LedgerChange.MustCreated().Data.MustReviewableRequest()
	histReviewableReq, err := c.convertReviewableRequest(&reviewableRequest, lc.LedgerCloseTime)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request", logan.F{
			"request":         reviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Insert(*histReviewableReq)
	if err != nil {
		return errors.Wrap(err, "failed to insert reviewable request", logan.F{
			"request":         histReviewableReq,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	op := lc.Operation.Body
	if op.Type == xdr.OperationTypeCreateKycRecoveryRequest {
		account := op.MustCreateKycRecoveryRequestOp().TargetAccount
		err = c.accounts.SetKYCRecoveryStatus(account.Address(), int(regources.KYCRecoveryStatusPending))
		if err != nil {
			return errors.Wrap(err, "failed to update status for account on create", logan.F{
				"request":         histReviewableReq,
				"ledger_sequence": lc.LedgerSeq,
			})
		}
	}

	return nil
}

// Updated - handles update of the request due to approve or reject op
func (c *reviewableRequestHandler) Updated(lc ledgerChange) error {
	reviewableRequest := lc.LedgerChange.MustUpdated().Data.MustReviewableRequest()
	histReviewableRequest, err := c.convertReviewableRequest(&reviewableRequest, lc.LedgerCloseTime)
	if err != nil {
		return errors.Wrap(err, "failed to convert reviewable request", logan.F{
			"request":         reviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	err = c.storage.Update(*histReviewableRequest)
	if err != nil {
		return errors.Wrap(err, "failed to update reviewable request", logan.F{
			"request":         histReviewableRequest,
			"ledger_sequence": lc.LedgerSeq,
		})
	}

	op := lc.Operation.Body
	if op.Type == xdr.OperationTypeReviewRequest {
		if op.MustReviewRequestOp().RequestDetails.RequestType == xdr.ReviewableRequestTypeKycRecovery {
			account := reviewableRequest.Body.MustKycRecoveryRequest().TargetAccount
			err = c.accounts.SetKYCRecoveryStatus(account.Address(), int(regources.KYCRecoveryStatusRejected))
			if err != nil {
				return errors.Wrap(err, "failed to update account status on update", logan.F{
					"request":         histReviewableRequest,
					"ledger_sequence": lc.LedgerSeq,
				})
			}
		}
	}

	return nil
}

func (c *reviewableRequestHandler) Stated(lc ledgerChange) error {
	request := lc.LedgerChange.MustState().Data.MustReviewableRequest()
	op := lc.Operation.Body
	switch op.Type {
	case xdr.OperationTypeReviewRequest:
		switch op.MustReviewRequestOp().RequestDetails.RequestType {
		case xdr.ReviewableRequestTypeKycRecovery:
			kycRec := request.Body.MustKycRecoveryRequest()
			address := kycRec.TargetAccount.Address()
			var err error
			switch op.ReviewRequestOp.Action {
			case xdr.ReviewRequestOpActionApprove:
				if lc.OperationResult.MustReviewRequestResult().MustSuccess().Fulfilled {
					err = c.accounts.SetKYCRecoveryStatus(address, int(regources.KYCRecoveryStatusNone))
				}
			case xdr.ReviewRequestOpActionPermanentReject:
				err = c.accounts.SetKYCRecoveryStatus(address, int(regources.KYCRecoveryStatusPermanentlyRejected))
			}

			if err != nil {
				return errors.Wrap(err, "failed to update account status on remove", logan.F{
					"request":         request,
					"ledger_sequence": lc.LedgerSeq,
				})
			}
		}
	}

	return nil
}

// Removed - handles removal of request from core due to approval, cancellation or permanent reject
func (c *reviewableRequestHandler) Removed(lc ledgerChange) error {
	// The request is deleted in 3 cases:
	// 1. Due to approve via reviewRequestOp
	// 2. Due to permanentReject via reviewRequestOp
	// 3. Due to cancel via specific operation
	op := lc.Operation.Body
	switch op.Type {
	case xdr.OperationTypeReviewRequest:
		return c.removedOnReview(lc)
	case xdr.OperationTypeManageAsset:
		return c.handleRemoveOnManageAsset(lc)
	case xdr.OperationTypeManageSale:
		return c.handleRemoveOnManageSale(lc)
	// auto review is handled by each operation separately
	case xdr.OperationTypeCreateIssuanceRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateIssuanceRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreatePreissuanceRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreatePreIssuanceRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateManageLimitsRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateManageLimitsRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateSaleRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateSaleCreationRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateAmlAlert:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateAmlAlertRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateChangeRoleRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateChangeRoleRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCreateAtomicSwapAskRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustCreateAtomicSwapAskRequestResult().MustSuccess().Fulfilled)
	case xdr.OperationTypeCheckSaleState:
		// if check sale state was successful, all the requests created by it were fulfilled
		return c.handleRemoveOnCreationOp(lc, true)
	case xdr.OperationTypeCancelSaleRequest:
		return c.cancel(lc)
	case xdr.OperationTypeManageCreatePollRequest:
		return c.handleRemoveOnManageCreatePollRequest(lc)
	case xdr.OperationTypeCancelChangeRoleRequest:
		return c.cancel(lc)
	case xdr.OperationTypeInitiateKycRecovery:
		return c.handleInitiateKycRecovery(lc)
	case xdr.OperationTypeCreateKycRecoveryRequest:
		{
			err := c.handleRemoveOnCreationOp(lc,
				lc.OperationResult.MustCreateKycRecoveryRequestResult().MustSuccess().Fulfilled)
			if err != nil {
				return err
			}
			account := op.MustCreateKycRecoveryRequestOp().TargetAccount
			return c.accounts.SetKYCRecoveryStatus(account.Address(), int(regources.KYCRecoveryStatusNone))
		}
	case xdr.OperationTypeCreateManageOfferRequest:
		return c.handleRemoveOnCreationOp(lc, true)
	case xdr.OperationTypeCreatePaymentRequest:
		return c.handleRemoveOnCreationOp(lc, true)
	case xdr.OperationTypeManageOffer:
		return c.handleRemovedOnManageOffer(lc)
	case xdr.OperationTypeCreateDataCreationRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.
				MustCreateDataCreationRequestResult().
				MustSuccess().Fulfilled)
	case xdr.OperationTypeCancelDataCreationRequest:
		return c.cancel(lc)
	case xdr.OperationTypeCreateDataUpdateRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.
				MustCreateDataUpdateRequestResult().
				MustSuccess().Fulfilled)
	case xdr.OperationTypeCancelDataUpdateRequest:
		return c.cancel(lc)
	case xdr.OperationTypeCreateDataRemoveRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.
				MustCreateDataUpdateRequestResult().
				MustSuccess().Fulfilled)
	case xdr.OperationTypeCancelDataRemoveRequest:
		return c.cancel(lc)
	case xdr.OperationTypeCancelDeferredPaymentCreationRequest:
		return c.cancel(lc)
	case xdr.OperationTypeCancelCloseDeferredPaymentRequest:
		return c.cancel(lc)
	case xdr.OperationTypeCreateDeferredPaymentCreationRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.
				MustCreateDeferredPaymentCreationRequestResult().
				MustSuccess().
				Fulfilled)
	case xdr.OperationTypeCreateCloseDeferredPaymentRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.
				MustCreateCloseDeferredPaymentRequestResult().
				MustSuccess().
				Fulfilled)
	case xdr.OperationTypeCreateRedemptionRequest:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.
				MustCreateRedemptionRequestResult().
				MustRedemptionResponse().
				Fulfilled)
	default: // safeguard for future updates
		return errors.From(errUnknownRemoveReason, logan.F{
			"op_type": op.Type.String(),
		})
	}
}

func (c *reviewableRequestHandler) handleInitiateKycRecovery(lc ledgerChange) error {
	id := uint64(lc.LedgerChange.MustRemoved().MustReviewableRequest().RequestId)
	return c.storage.PermanentReject(id, removeOnKYCRecoveryInit)
}

func (c *reviewableRequestHandler) handleRemoveOnCreationOp(lc ledgerChange, fulfilled bool) error {
	id := uint64(lc.LedgerChange.MustRemoved().MustReviewableRequest().RequestId)
	if !fulfilled {
		return errors.From(errors.New("unexpected state: request has been removed on creation op, but fulfilled is false"), logan.F{
			"id": uint64(id),
		})
	}

	return c.storage.Approve(id)
}

func (c *reviewableRequestHandler) handleRemovedOnManageOffer(lc ledgerChange) error {
	id := uint64(lc.LedgerChange.MustRemoved().MustReviewableRequest().RequestId)
	if lc.Operation.Body.MustManageOfferOp().OrderBookId == 0 {
		return errors.From(errors.New("unexpected state: request has been removed on manage offer op, "+
			"but order book is secondary"), logan.F{
			"id": uint64(id),
		})
	}
	return c.storage.Approve(id)
}

func (c *reviewableRequestHandler) handleRemoveOnManageCreatePollRequest(lc ledgerChange) error {
	data := lc.Operation.Body.MustManageCreatePollRequestOp().Data
	switch data.Action {
	case xdr.ManageCreatePollRequestActionCreate:
		return c.handleRemoveOnCreationOp(lc,
			lc.OperationResult.MustManageCreatePollRequestResult().
				MustSuccess().Details.MustResponse().Fulfilled)
	case xdr.ManageCreatePollRequestActionCancel:
		return c.cancel(lc)
	default:
		return errors.Wrap(errUnknownRemoveReason,
			"Unexpected manage create poll request action", logan.F{
				"action": data.Action.String(),
			})
	}
}

func (c *reviewableRequestHandler) handleRemoveOnManageAsset(lc ledgerChange) error {
	op := lc.Operation.Body.MustManageAssetOp()
	switch op.Request.Action {
	// must be handled by operation
	case xdr.ManageAssetActionCreateAssetCreationRequest,
		xdr.ManageAssetActionCreateAssetUpdateRequest,
		xdr.ManageAssetActionChangePreissuedAssetSigner,
		xdr.ManageAssetActionUpdateMaxIssuance:
		fulfilled := lc.OperationResult.MustManageAssetResult().MustSuccess().Fulfilled
		return c.handleRemoveOnCreationOp(lc, fulfilled)
	case xdr.ManageAssetActionCancelAssetRequest:
		return c.cancel(lc)
	default:
		return errors.From(errUnknownRemoveReason, logan.F{
			"manage_asset_action": op.Request.Action,
		})
	}
}

func (c *reviewableRequestHandler) handleRemoveOnManageSale(lc ledgerChange) error {
	op := lc.Operation.Body.MustManageSaleOp()
	switch op.Data.Action {
	// must be handled by operation
	case xdr.ManageSaleActionCreateUpdateDetailsRequest:
		fulfilled := lc.OperationResult.MustManageSaleResult().MustSuccess().Fulfilled
		return c.handleRemoveOnCreationOp(lc, fulfilled)
	case xdr.ManageSaleActionCancel:
		return c.cancel(lc)
	default:
		return errors.From(errUnknownRemoveReason, logan.F{
			"manage_sale_action": op.Data.Action,
		})
	}
}

func (c *reviewableRequestHandler) removedOnReview(lc ledgerChange) error {
	key := lc.LedgerChange.MustRemoved().MustReviewableRequest()
	op := lc.Operation.Body.MustReviewRequestOp()

	switch op.Action {
	case xdr.ReviewRequestOpActionApprove:
		err := c.storage.Approve(uint64(key.RequestId))
		if err != nil {
			return errors.Wrap(err, "Failed to delete reviewable request", logan.F{
				"ledger_entry_key": key,
			})
		}
	case xdr.ReviewRequestOpActionPermanentReject:
		err := c.storage.PermanentReject(uint64(key.RequestId), string(op.Reason))
		if err != nil {
			return errors.Wrap(err, "Failed to delete reviewable request", logan.F{
				"ledger_entry_key": key,
			})
		}
	default:
		return errors.From(
			errors.New("unknown action during handle of removed reviewable request on review operation"),
			logan.F{
				"action": op.Action.String(),
			})
	}

	return nil
}

func (c *reviewableRequestHandler) cancel(lc ledgerChange) error {
	key := lc.LedgerChange.MustRemoved().MustReviewableRequest()
	err := c.storage.Cancel(uint64(key.RequestId))
	if err != nil {
		return errors.Wrap(err, "failed to cancel reivewalbe request on remove by operation", logan.F{
			"op_type": lc.Operation.Body.Type,
		})
	}

	return nil
}

func (c *reviewableRequestHandler) convertReviewableRequest(request *xdr.ReviewableRequestEntry,
	ledgerCloseTime time.Time,
) (*history.ReviewableRequest, error) {

	var reference *string
	if request.Reference != nil {
		reference = new(string)
		*reference = utf8.Scrub(string(*request.Reference))
	}

	details, err := c.getReviewableRequestDetails(&request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reviewable request details")
	}

	state := history.ReviewableRequestStatePending
	if string(request.RejectReason) != "" {
		state = history.ReviewableRequestStateRejected
	}

	result := history.ReviewableRequest{
		TotalOrderID: db2.TotalOrderID{
			ID: int64(request.RequestId),
		},
		Requestor:    request.Requestor.Address(),
		Reviewer:     request.Reviewer.Address(),
		Reference:    reference,
		RejectReason: string(request.RejectReason),
		RequestType:  request.Body.Type,
		RequestState: state,
		Hash:         hex.EncodeToString(request.Hash[:]),
		Details:      details,
		CreatedAt:    unixToTime(int64(request.CreatedAt)),
		UpdatedAt:    ledgerCloseTime,
	}

	tasksExt := request.Tasks
	result.AllTasks = uint32(tasksExt.AllTasks)
	result.PendingTasks = uint32(tasksExt.PendingTasks)

	externalDetails := make([]regources.Details, 0, len(tasksExt.ExternalDetails))
	for _, item := range tasksExt.ExternalDetails {
		externalDetails = append(externalDetails, internal.MarshalCustomDetails(item))
	}

	// we use key "data" for compatibility with db2.Details (Deprecated)
	// the value for the key "data" is a slice of map[string]interface{}
	var resultDetails = map[string]interface{}{
		"data": externalDetails,
	}
	result.ExternalDetails, err = json.Marshal(resultDetails)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *reviewableRequestHandler) getAssetCreation(request *xdr.AssetCreationRequest) *history.CreateAssetRequest {
	return &history.CreateAssetRequest{
		Asset:                  string(request.Code),
		Type:                   uint64(request.Type),
		Policies:               int32(request.Policies),
		PreIssuedAssetSigner:   request.PreissuedAssetSigner.Address(),
		MaxIssuanceAmount:      regources.Amount(request.MaxIssuanceAmount),
		InitialPreissuedAmount: regources.Amount(request.InitialPreissuedAmount),
		CreatorDetails:         internal.MarshalCustomDetails(request.CreatorDetails),
		TrailingDigitsCount:    uint32(request.TrailingDigitsCount),
	}
}

func (c *reviewableRequestHandler) getAssetUpdate(request *xdr.AssetUpdateRequest) *history.UpdateAssetRequest {
	return &history.UpdateAssetRequest{
		Asset:          string(request.Code),
		Policies:       int32(request.Policies),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getPreIssuanceRequest(request *xdr.PreIssuanceRequest,
) (*history.CreatePreIssuanceRequest, error) {

	signature, err := xdr.MarshalBase64(request.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal signature")
	}

	return &history.CreatePreIssuanceRequest{
		Asset:          string(request.Asset),
		Amount:         regources.Amount(request.Amount),
		Signature:      signature,
		Reference:      string(request.Reference),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}, nil
}

func (c *reviewableRequestHandler) getIssuanceRequest(request *xdr.IssuanceRequest) *history.CreateIssuanceRequest {
	return &history.CreateIssuanceRequest{
		Asset:          string(request.Asset),
		Amount:         regources.Amount(request.Amount),
		Receiver:       request.Receiver.AsString(),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getWithdrawalRequest(request *xdr.WithdrawalRequest) *history.CreateWithdrawalRequest {
	histBalance := c.balances.MustBalance(request.Balance)

	return &history.CreateWithdrawalRequest{
		Asset:     histBalance.AssetCode,
		BalanceID: request.Balance.AsString(),
		Amount:    regources.Amount(request.Amount),
		Fee: regources.Fee{
			Fixed:             regources.Amount(request.Fee.Fixed),
			CalculatedPercent: regources.Amount(request.Fee.Percent),
		},
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAmlAlertRequest(request *xdr.AmlAlertRequest) *history.CreateAmlAlertRequest {
	return &history.CreateAmlAlertRequest{
		BalanceID:      request.BalanceId.AsString(),
		Amount:         regources.Amount(request.Amount),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getSaleRequest(request *xdr.SaleCreationRequest) *history.CreateSaleRequest {
	quoteAssets := make([]regources.AssetPrice, 0, len(request.QuoteAssets))
	for i := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Price: regources.Amount(int64(request.QuoteAssets[i].Price)),
			Asset: string(request.QuoteAssets[i].QuoteAsset),
		})
	}

	saleType := request.SaleTypeExt.SaleType
	accessDefinitionType := getSaleAccessDefinitionType(request)

	return &history.CreateSaleRequest{
		BaseAsset:            string(request.BaseAsset),
		DefaultQuoteAsset:    string(request.DefaultQuoteAsset),
		StartTime:            unixToTime(int64(request.StartTime)),
		EndTime:              unixToTime(int64(request.EndTime)),
		SoftCap:              regources.Amount(request.SoftCap),
		HardCap:              regources.Amount(request.HardCap),
		CreatorDetails:       internal.MarshalCustomDetails(request.CreatorDetails),
		QuoteAssets:          quoteAssets,
		SaleType:             saleType,
		BaseAssetForHardCap:  regources.Amount(request.RequiredBaseAssetForHardCap),
		AccessDefinitionType: accessDefinitionType,
	}
}

func getSaleAccessDefinitionType(request *xdr.SaleCreationRequest) regources.SaleAccessDefinitionType {
	if request.Ext.SaleRules == nil {
		return regources.SaleAccessDefinitionTypeNone
	}

	for _, rule := range *request.Ext.SaleRules {
		if rule.AccountId != nil {
			continue
		}
		if rule.Forbids {
			return regources.SaleAccessDefinitionTypeWhitelist
		} else {
			return regources.SaleAccessDefinitionTypeBlacklist
		}
	}

	return regources.SaleAccessDefinitionTypeNone
}

func (c *reviewableRequestHandler) getLimitsUpdateRequest(request *xdr.LimitsUpdateRequest,
) *history.UpdateLimitsRequest {

	return &history.UpdateLimitsRequest{
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getChangeRoleRequest(request *xdr.ChangeRoleRequest,
) *history.ChangeRoleRequest {
	return &history.ChangeRoleRequest{
		DestinationAccount: request.DestinationAccount.Address(),
		AccountRoleToSet:   uint64(request.AccountRoleToSet),
		CreatorDetails:     internal.MarshalCustomDetails(request.CreatorDetails),
		SequenceNumber:     uint32(request.SequenceNumber),
	}
}

func (c *reviewableRequestHandler) getUpdateSaleDetailsRequest(
	request *xdr.UpdateSaleDetailsRequest) *history.UpdateSaleDetailsRequest {
	return &history.UpdateSaleDetailsRequest{
		SaleID:         uint64(request.SaleId),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getCreatePollRequest(
	request *xdr.CreatePollRequest) *history.CreatePollRequest {
	return &history.CreatePollRequest{
		NumberOfChoices:          uint32(request.NumberOfChoices),
		VoteConfirmationRequired: request.VoteConfirmationRequired,
		ResultProviderID:         request.ResultProviderId.Address(),
		PermissionType:           uint32(request.PermissionType),
		PollData:                 request.Data,
		StartTime:                unixToTime(int64(request.StartTime)),
		EndTime:                  unixToTime(int64(request.EndTime)),
		CreatorDetails:           internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getAtomicSwapAskCreationRequest(request *xdr.CreateAtomicSwapAskRequest,
) *history.CreateAtomicSwapAskRequest {
	quoteAssets := make([]regources.AssetPrice, 0, len(request.QuoteAssets))
	for _, quoteAsset := range request.QuoteAssets {
		quoteAssets = append(quoteAssets, regources.AssetPrice{
			Asset: string(quoteAsset.QuoteAsset),
			Price: regources.Amount(quoteAsset.Price),
		})
	}

	return &history.CreateAtomicSwapAskRequest{
		BaseBalance:    request.BaseBalance.AsString(),
		BaseAmount:     regources.Amount(request.Amount),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
		QuoteAssets:    quoteAssets,
	}
}

func (c *reviewableRequestHandler) getAtomicSwapBidRequest(request *xdr.CreateAtomicSwapBidRequest,
) *history.CreateAtomicSwapBidRequest {
	return &history.CreateAtomicSwapBidRequest{
		AskID:          uint64(request.AskId),
		BaseAmount:     regources.Amount(request.BaseAmount),
		QuoteAsset:     string(request.QuoteAsset),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getKYCRecovery(request *xdr.KycRecoveryRequest,
) *history.KYCRecoveryRequest {
	signersData := make([]history.UpdateSignerDetails, 0, len(request.SignersData))
	for _, signer := range request.SignersData {
		signersData = append(signersData, history.UpdateSignerDetails{
			Details:  internal.MarshalCustomDetails(signer.Details),
			RoleID:   uint64(signer.RoleId),
			Identity: uint32(signer.Identity),
			Weight:   uint32(signer.Weight),
		})
	}

	return &history.KYCRecoveryRequest{
		TargetAccount:  request.TargetAccount.Address(),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
		SequenceNumber: uint32(request.SequenceNumber),
		SignersData:    signersData,
	}
}

func (c *reviewableRequestHandler) getManageOfferRequest(request *xdr.ManageOfferRequest,
) *history.ManageOfferRequest {
	manageOfferOp := request.Op

	creatorDetails := regources.Details("{}")
	switch request.Ext.V {
	case xdr.LedgerVersionMovementRequestsDetails:
		creatorDetails = internal.MarshalCustomDetails(request.Ext.MustCreatorDetails())
	case xdr.LedgerVersionEmptyVersion:
	default:
		panic(errors.From(errors.New("unexpected version of manage offer request"), logan.F{
			"ledger_version": request.Ext.V,
		}))
	}

	return &history.ManageOfferRequest{
		CreatorDetails: creatorDetails,
		OfferID:        int64(manageOfferOp.OfferId),
		OrderBookID:    int64(manageOfferOp.OrderBookId),
		Amount:         regources.Amount(manageOfferOp.Amount),
		Price:          regources.Amount(manageOfferOp.Price),
		IsBuy:          manageOfferOp.IsBuy,
		Fee: regources.Fee{
			CalculatedPercent: regources.Amount(manageOfferOp.Fee),
		},
	}
}

func (c *reviewableRequestHandler) getCreatePaymentRequest(request *xdr.CreatePaymentRequest,
) *history.CreatePaymentRequest {
	paymentOp := request.PaymentOp

	creatorDetails := regources.Details("{}")
	switch request.Ext.V {
	case xdr.LedgerVersionMovementRequestsDetails:
		creatorDetails = internal.MarshalCustomDetails(request.Ext.MustCreatorDetails())
	case xdr.LedgerVersionEmptyVersion:
	default:
		panic(errors.From(errors.New("unexpected version of payment request"), logan.F{
			"ledger_version": request.Ext.V,
		}))
	}

	return &history.CreatePaymentRequest{
		CreatorDetails:          creatorDetails,
		BalanceFrom:             paymentOp.SourceBalanceId.AsString(),
		Amount:                  regources.Amount(paymentOp.Amount),
		SourceFee:               internal.FeeFromXdr(paymentOp.FeeData.SourceFee),
		DestinationFee:          internal.FeeFromXdr(paymentOp.FeeData.DestinationFee),
		SourcePayForDestination: paymentOp.FeeData.SourcePaysForDest,
		Subject:                 string(paymentOp.Subject),
		Reference:               utf8.Scrub(string(paymentOp.Reference)),
	}
}

func (c *reviewableRequestHandler) getRedemption(request *xdr.RedemptionRequest) *history.RedemptionRequest {
	return &history.RedemptionRequest{
		SourceBalanceID:      request.SourceBalanceId.AsString(),
		DestinationAccountID: request.Destination.Address(),
		Amount:               regources.Amount(request.Amount),
		CreatorDetails:       regources.Details(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getDataCreationRequest(request *xdr.DataCreationRequest) *history.DataCreationRequest {
	return &history.DataCreationRequest{
		SecurityType:   uint64(request.Type),
		SequenceNumber: uint32(request.SequenceNumber),
		Owner:          request.Owner.Address(),
		Value:          internal.MarshalCustomDetails(request.Value),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getDataUpdateRequest(request *xdr.DataUpdateRequest) *history.DataUpdateRequest {
	return &history.DataUpdateRequest{
		SequenceNumber: uint32(request.SequenceNumber),
		DataID:         uint64(request.Id),
		Value:          internal.MarshalCustomDetails(request.Value),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getDataRemoveRequest(request *xdr.DataRemoveRequest) *history.DataRemoveRequest {
	return &history.DataRemoveRequest{
		SequenceNumber: uint32(request.SequenceNumber),
		DataID:         uint64(request.Id),
		CreatorDetails: internal.MarshalCustomDetails(request.CreatorDetails),
	}
}

func (c *reviewableRequestHandler) getCreateDeferredPayment(request *xdr.CreateDeferredPaymentRequest) *history.CreateDeferredPayment {
	return &history.CreateDeferredPayment{
		SourceBalance:      request.SourceBalance.AsString(),
		DestinationAccount: request.Destination.Address(),
		Amount:             regources.Amount(request.Amount),
		Details:            internal.MarshalCustomDetails(request.CreatorDetails),
		SequenceNumber:     uint32(request.SequenceNumber),
	}
}

func (c *reviewableRequestHandler) getCloseDeferredPayment(request *xdr.CloseDeferredPaymentRequest) *history.CloseDeferredPayment {
	result := &history.CloseDeferredPayment{
		SequenceNumber:    uint32(request.SequenceNumber),
		DeferredPaymentID: uint64(request.DeferredPaymentId),
		Destination: history.CloseDeferredPaymentDestination{
			Type: request.Destination.Type,
		},
		Amount:  regources.Amount(request.Amount),
		Details: internal.MarshalCustomDetails(request.CreatorDetails),
	}

	switch request.Destination.Type {
	case xdr.CloseDeferredPaymentDestinationTypeAccount:
		tmp := request.Destination.AccountId.Address()
		result.Destination.Account = &tmp
	case xdr.CloseDeferredPaymentDestinationTypeBalance:
		tmp := request.Destination.BalanceId.AsString()
		result.Destination.Balance = &tmp
	}

	return result
}
func (c *reviewableRequestHandler) getReviewableRequestDetails(
	body *xdr.ReviewableRequestEntryBody,
) (history.ReviewableRequestDetails, error) {

	details := history.ReviewableRequestDetails{
		Type: body.Type,
	}

	var err error
	switch body.Type {
	case xdr.ReviewableRequestTypeCreateAsset:
		details.CreateAsset = c.getAssetCreation(body.AssetCreationRequest)
	case xdr.ReviewableRequestTypeUpdateAsset:
		details.UpdateAsset = c.getAssetUpdate(body.AssetUpdateRequest)
	case xdr.ReviewableRequestTypeCreateIssuance:
		details.CreateIssuance = c.getIssuanceRequest(body.IssuanceRequest)
	case xdr.ReviewableRequestTypeCreatePreIssuance:
		details.CreatePreIssuance, err = c.getPreIssuanceRequest(body.PreIssuanceRequest)
		if err != nil {
			return details, errors.Wrap(err, "failed to get pre issuance request")
		}
	case xdr.ReviewableRequestTypeCreateWithdraw:
		details.CreateWithdraw = c.getWithdrawalRequest(body.WithdrawalRequest)
	case xdr.ReviewableRequestTypeCreateSale:
		details.CreateSale = c.getSaleRequest(body.SaleCreationRequest)
	case xdr.ReviewableRequestTypeUpdateLimits:
		details.UpdateLimits = c.getLimitsUpdateRequest(body.LimitsUpdateRequest)
	case xdr.ReviewableRequestTypeCreateAmlAlert:
		details.CreateAmlAlert = c.getAmlAlertRequest(body.AmlAlertRequest)
	case xdr.ReviewableRequestTypeChangeRole:
		details.ChangeRole = c.getChangeRoleRequest(body.ChangeRoleRequest)
	case xdr.ReviewableRequestTypeUpdateSaleDetails:
		details.UpdateSaleDetails = c.getUpdateSaleDetailsRequest(body.UpdateSaleDetailsRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwapAsk:
		details.CreateAtomicSwapAsk = c.getAtomicSwapAskCreationRequest(body.CreateAtomicSwapAskRequest)
	case xdr.ReviewableRequestTypeCreateAtomicSwapBid:
		details.CreateAtomicSwapBid = c.getAtomicSwapBidRequest(body.CreateAtomicSwapBidRequest)
	case xdr.ReviewableRequestTypeCreatePoll:
		details.CreatePoll = c.getCreatePollRequest(body.CreatePollRequest)
	case xdr.ReviewableRequestTypeKycRecovery:
		details.KYCRecovery = c.getKYCRecovery(body.KycRecoveryRequest)
	case xdr.ReviewableRequestTypeManageOffer:
		details.ManageOffer = c.getManageOfferRequest(body.ManageOfferRequest)
	case xdr.ReviewableRequestTypeCreatePayment:
		details.CreatePayment = c.getCreatePaymentRequest(body.CreatePaymentRequest)
	case xdr.ReviewableRequestTypePerformRedemption:
		details.Redemption = c.getRedemption(body.RedemptionRequest)
	case xdr.ReviewableRequestTypeDataCreation:
		details.DataCreation = c.getDataCreationRequest(body.DataCreationRequest)
	case xdr.ReviewableRequestTypeDataUpdate:
		details.DataUpdate = c.getDataUpdateRequest(body.DataUpdateRequest)
	case xdr.ReviewableRequestTypeDataRemove:
		details.DataRemove = c.getDataRemoveRequest(body.DataRemoveRequest)
	case xdr.ReviewableRequestTypeCreateDeferredPayment:
		details.CreateDeferredPayment = c.getCreateDeferredPayment(body.CreateDeferredPaymentRequest)
	case xdr.ReviewableRequestTypeCloseDeferredPayment:
		details.CloseDeferredPayment = c.getCloseDeferredPayment(body.CloseDeferredPaymentRequest)
	default:
		return details, errors.From(errors.New("unexpected reviewable request type"), map[string]interface{}{
			"request_type": body.Type.String(),
		})
	}

	return details, nil
}
