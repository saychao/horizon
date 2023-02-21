package history2

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"

	"database/sql/driver"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

//ReviewableRequestDetails - stores in union switch details of the reviewable requests
type ReviewableRequestDetails struct {
	Type                  xdr.ReviewableRequestType   `json:"type"`
	CreateAsset           *CreateAssetRequest         `json:"create_asset,omitempty"`
	UpdateAsset           *UpdateAssetRequest         `json:"update_asset,omitempty"`
	CreatePreIssuance     *CreatePreIssuanceRequest   `json:"create_pre_issuance,omitempty"`
	CreateIssuance        *CreateIssuanceRequest      `json:"create_issuance,omitempty"`
	CreateWithdraw        *CreateWithdrawalRequest    `json:"create_withdraw,omitempty"`
	CreateSale            *CreateSaleRequest          `json:"create_sale,omitempty"`
	UpdateLimits          *UpdateLimitsRequest        `json:"update_limits,omitempty"`
	CreateAmlAlert        *CreateAmlAlertRequest      `json:"create_aml_alert,omitempty"`
	ChangeRole            *ChangeRoleRequest          `json:"change_role,omitempty"`
	UpdateSaleDetails     *UpdateSaleDetailsRequest   `json:"update_sale_details,omitempty"`
	CreateAtomicSwapAsk   *CreateAtomicSwapAskRequest `json:"create_atomic_swap_ask,omitempty"`
	CreateAtomicSwapBid   *CreateAtomicSwapBidRequest `json:"create_atomic_swap_bid,omitempty"`
	CreatePoll            *CreatePollRequest          `json:"create_poll,omitempty"`
	KYCRecovery           *KYCRecoveryRequest         `json:"kyc_recovery,omitempty"`
	ManageOffer           *ManageOfferRequest         `json:"manage_offer,omitempty"`
	CreatePayment         *CreatePaymentRequest       `json:"create_payment,omitempty"`
	Redemption            *RedemptionRequest          `json:"redemption,omitempty"`
	DataCreation          *DataCreationRequest        `json:"data_creation,omitempty"`
	DataUpdate            *DataUpdateRequest          `json:"data_update,omitempty"`
	DataRemove            *DataRemoveRequest          `json:"data_remove,omitempty"`
	CreateDeferredPayment *CreateDeferredPayment      `json:"create_deferred_payment,omitempty"`
	CloseDeferredPayment  *CloseDeferredPayment       `json:"close_deferred_payment,omitempty"`
}

//Value - implements db driver method for auto marshal
func (r ReviewableRequestDetails) Value() (driver.Value, error) {
	result, err := pgdb.JSONValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *ReviewableRequestDetails) Scan(src interface{}) error {
	err := pgdb.JSONScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details")
	}

	return nil
}

//CreateAssetRequest - asset creation request details
type CreateAssetRequest struct {
	Asset                  string            `json:"asset"`
	Type                   uint64            `json:"type"`
	Policies               int32             `json:"policies"`
	PreIssuedAssetSigner   string            `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      regources.Amount  `json:"max_issuance_amount"`
	InitialPreissuedAmount regources.Amount  `json:"initial_preissued_amount"`
	CreatorDetails         regources.Details `json:"details"`
	TrailingDigitsCount    uint32            `json:"trailing_digits_count"`
}

//UpdateAssetRequest - asset update request details
type UpdateAssetRequest struct {
	Asset          string            `json:"asset"`
	Policies       int32             `json:"policies"`
	CreatorDetails regources.Details `json:"details"`
}

//CreatePreIssuanceRequest - request details
type CreatePreIssuanceRequest struct {
	Asset          string            `json:"asset"`
	Amount         regources.Amount  `json:"amount"`
	Signature      string            `json:"signature"`
	Reference      string            `json:"reference"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateIssuanceRequest - request details
type CreateIssuanceRequest struct {
	Asset          string            `json:"asset"`
	Amount         regources.Amount  `json:"amount"`
	Receiver       string            `json:"receiver"`
	CreatorDetails regources.Details `json:"external_details"`
}

//CreateWithdrawalRequest - request details
type CreateWithdrawalRequest struct {
	Asset          string            `json:"asset"`
	BalanceID      string            `json:"balance_id"`
	Amount         regources.Amount  `json:"amount"`
	Fee            regources.Fee     `json:"fee"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateSaleRequest - request details
type CreateSaleRequest struct {
	BaseAsset            string                             `json:"base_asset"`
	DefaultQuoteAsset    string                             `json:"quote_asset"`
	StartTime            time.Time                          `json:"start_time"`
	EndTime              time.Time                          `json:"end_time"`
	SoftCap              regources.Amount                   `json:"soft_cap"`
	HardCap              regources.Amount                   `json:"hard_cap"`
	CreatorDetails       regources.Details                  `json:"creator_details"`
	QuoteAssets          []regources.AssetPrice             `json:"quote_assets"`
	SaleType             xdr.SaleType                       `json:"sale_type"`
	BaseAssetForHardCap  regources.Amount                   `json:"base_asset_for_hard_cap"`
	AccessDefinitionType regources.SaleAccessDefinitionType `json:"access_definition_type"`
}

//UpdateLimitsRequest - request details
type UpdateLimitsRequest struct {
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateAmlAlertRequest - request details
type CreateAmlAlertRequest struct {
	BalanceID      string            `json:"balance_id"`
	Amount         regources.Amount  `json:"amount"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//ChangeRoleRequest - request details
type ChangeRoleRequest struct {
	DestinationAccount string            `json:"destination_account"`
	AccountRoleToSet   uint64            `json:"account_role_to_set"`
	SequenceNumber     uint32            `json:"sequence_number"`
	CreatorDetails     regources.Details `json:"creator_details"`
}

//UpdateSaleDetailsRequest - request details
type UpdateSaleDetailsRequest struct {
	SaleID         uint64            `json:"sale_id"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateAtomicSwapAskRequest - request details
type CreateAtomicSwapAskRequest struct {
	BaseBalance    string                 `json:"base_balance"`
	BaseAmount     regources.Amount       `json:"base_amount"`
	CreatorDetails regources.Details      `json:"creator_details"`
	QuoteAssets    []regources.AssetPrice `json:"quote_assets"`
}

//CreateAtomicSwapAskRequest - request details
type CreateAtomicSwapBidRequest struct {
	AskID          uint64            `json:"bid_id"`
	BaseAmount     regources.Amount  `json:"base_amount"`
	QuoteAsset     string            `json:"quote_asset"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreatePollRequest - request details
type CreatePollRequest struct {
	PermissionType           uint32            `json:"permission_type"`
	NumberOfChoices          uint32            `json:"number_of_choices"`
	PollData                 xdr.PollData      `json:"poll_data"`
	CreatorDetails           regources.Details `json:"creator_details"`
	VoteConfirmationRequired bool              `json:"vote_confirmation_required"`
	ResultProviderID         string            `json:"result_provider_id"`
	StartTime                time.Time         `json:"start_time"`
	EndTime                  time.Time         `json:"end_time"`
}

type KYCRecoveryRequest struct {
	TargetAccount  string                `json:"target_account"`
	SignersData    []UpdateSignerDetails `json:"signers_data"`
	CreatorDetails regources.Details     `json:"creator_details"`
	SequenceNumber uint32                `json:"sequence_number"`
}

type ManageOfferRequest struct {
	CreatorDetails regources.Details `json:"creator_details"`
	OfferID        int64             `json:"offer_id,omitempty"`
	OrderBookID    int64             `json:"order_book_id"`
	Amount         regources.Amount  `json:"base_amount"`
	Price          regources.Amount  `json:"price"`
	IsBuy          bool              `json:"is_buy"`
	Fee            regources.Fee     `json:"fee"`
}

type CreatePaymentRequest struct {
	CreatorDetails          regources.Details `json:"creator_details"`
	BalanceFrom             string            `json:"balance_from"`
	Amount                  regources.Amount  `json:"amount"`
	SourceFee               regources.Fee     `json:"source_fee"`
	DestinationFee          regources.Fee     `json:"destination_fee"`
	SourcePayForDestination bool              `json:"source_pay_for_destination"`
	Subject                 string            `json:"subject"`
	Reference               string            `json:"reference"`
}

type RedemptionRequest struct {
	SourceBalanceID      string            `json:"source_balance_id"`
	DestinationAccountID string            `json:"destination_account_id"`
	Amount               regources.Amount  `json:"amount"`
	CreatorDetails       regources.Details `json:"creator_details"`
}

type DataCreationRequest struct {
	SecurityType   uint64            `json:"security_type"`
	SequenceNumber uint32            `json:"sequence_number"`
	Owner          string            `json:"owner"`
	Value          regources.Details `json:"value"`
	CreatorDetails regources.Details `json:"creator_details"`
}

type DataUpdateRequest struct {
	SequenceNumber uint32            `json:"sequence_number"`
	Value          regources.Details `json:"value"`
	DataID         uint64            `json:"data_id"`
	CreatorDetails regources.Details `json:"creator_details"`
}

type DataRemoveRequest struct {
	SequenceNumber uint32            `json:"sequence_number"`
	DataID         uint64            `json:"data_id"`
	CreatorDetails regources.Details `json:"creator_details"`
}

type CreateDeferredPayment struct {
	SourceBalance      string            `json:"source_balance"`
	DestinationAccount string            `json:"destination_account"`
	Amount             regources.Amount  `json:"amount"`
	Details            regources.Details `json:"details"`
	SequenceNumber     uint32            `json:"sequence_number"`
}

type CloseDeferredPayment struct {
	SequenceNumber    uint32                          `json:"sequence_number"`
	Destination       CloseDeferredPaymentDestination `json:"destination"`
	DeferredPaymentID uint64                          `json:"deferred_payment_id"`
	Amount            regources.Amount                `json:"amount"`
	Details           regources.Details               `json:"details"`
}

type CloseDeferredPaymentDestination struct {
	Type    xdr.CloseDeferredPaymentDestinationType `json:"type"`
	Balance *string                                 `json:"balance,omitempty"`
	Account *string                                 `json:"account,omitempty"`
}
