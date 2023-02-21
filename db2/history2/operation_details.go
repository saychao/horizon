package history2

import (
	"database/sql/driver"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

// OperationDetails - stores details of the operation performed in union switch form.
// Only one value must be selected at a type
type OperationDetails struct {
	//NOTE: omitempty MUST be specified for each switch value
	Type                                 xdr.OperationType                     `json:"type"`
	CreateAccount                        *CreateAccountDetails                 `json:"create_account,omitempty"`
	ManageAccountRule                    *ManageAccountRuleDetails             `json:"manage_account_rule,omitempty"`
	ManageAccountRole                    *ManageAccountRoleDetails             `json:"manage_account_role,omitempty"`
	ManageAccountSpecificRule            *ManageAccountSpecificRuleDetails     `json:"manage_account_specific_rule,omitempty"`
	ManageSigner                         *ManageSignerDetails                  `json:"manage_signer,omitempty"`
	ManageSignerRule                     *ManageSignerRuleDetails              `json:"manage_signer_rule,omitempty"`
	ManageSignerRole                     *ManageSignerRoleDetails              `json:"manage_signer_role,omitempty"`
	ManageBalance                        *ManageBalanceDetails                 `json:"manage_balance,omitempty"`
	ManageKeyValue                       *ManageKeyValueDetails                `json:"manage_key_value,omitempty"`
	ManageAsset                          *ManageAssetDetails                   `json:"manage_asset,omitempty"`
	ManageAssetPair                      *ManageAssetPairDetails               `json:"manage_asset_pair,omitempty"`
	ManageLimits                         *ManageLimitsDetails                  `json:"manage_limits,omitempty"`
	ManageOffer                          *ManageOfferDetails                   `json:"manage_offer,omitempty"`
	ManageExternalSystemPool             *ManageExternalSystemPoolDetails      `json:"manage_external_system_pool,omitempty"`
	BindExternalSystemAccount            *BindExternalSystemAccountDetails     `json:"bind_external_system_account,omitempty"`
	Payment                              *PaymentDetails                       `json:"payment,omitempty"`
	Payout                               *PayoutDetails                        `json:"payout,omitempty"`
	SetFee                               *SetFeeDetails                        `json:"set_fee,omitempty"`
	CancelAtomicSwapAsk                  *CancelAtomicSwapAskDetails           `json:"cancel_atomic_swap_ask,omitempty"`
	CheckSaleState                       *CheckSaleStateDetails                `json:"check_sale_state,omitempty"`
	CreateChangeRoleRequest              *CreateChangeRoleRequestDetails       `json:"create_change_role_request,omitempty"`
	CreatePreIssuanceRequest             *CreatePreIssuanceRequestDetails      `json:"create_pre_issuance_request,omitempty"`
	CreateIssuanceRequest                *CreateIssuanceRequestDetails         `json:"create_issuance_request,omitempty"`
	CreateManageLimitsRequest            *CreateManageLimitsRequestDetails     `json:"create_manage_limits_request,omitempty"`
	CreateWithdrawRequest                *CreateWithdrawRequestDetails         `json:"create_withdraw_request,omitempty"`
	CreateAMLAlertRequest                *CreateAMLAlertRequestDetails         `json:"create_aml_alert_request,omitempty"`
	CreateSaleRequest                    *CreateSaleRequestDetails             `json:"create_sale_request,omitempty"`
	CreateAtomicSwapAskRequest           *CreateAtomicSwapAskRequestDetails    `json:"create_atomic_swap_ask_request,omitempty"`
	CreateAtomicSwapBidRequest           *CreateAtomicSwapBidRequestDetails    `json:"create_atomic_swap_bid_request,omitempty"`
	ReviewRequest                        *ReviewRequestDetails                 `json:"review_request,omitempty"`
	ManageSale                           *ManageSaleDetails                    `json:"manage_sale,omitempty"`
	License                              *LicenseDetails                       `json:"license,omitempty"`
	Stamp                                *StampDetails                         `json:"stamp,omitempty"`
	ManageCreatePollRequest              *ManageCreatePollRequestDetails       `json:"manage_create_poll_request,omitempty"`
	ManagePoll                           *ManagePollDetails                    `json:"manage_poll,omitempty"`
	ManageVote                           *ManageVoteDetails                    `json:"manage_vote,omitempty"`
	RemoveAssetPair                      *RemoveAssetPairDetails               `json:"remove_asset_pair,omitempty"`
	InitiateKYCRecovery                  *InitiateKYCRecoveryDetails           `json:"initiate_kyc_recovery,omitempty"`
	CreateKYCRecoveryRequest             *CreateKYCRecoveryRequestDetails      `json:"create_kyc_recovery_request,omitempty"`
	CreateManageOfferRequest             *CreateManageOfferRequestDetails      `json:"create_manage_offer_request,omitempty"`
	CreatePaymentRequest                 *CreatePaymentRequestDetails          `json:"create_payment_request,omitempty"`
	RemoveAsset                          *RemoveAssetDetails                   `json:"remove_asset,omitempty"`
	OpenSwap                             *OpenSwapDetails                      `json:"open_swap,omitempty"`
	CloseSwap                            *CloseSwapDetails                     `json:"close_swap,omitempty"`
	Redemption                           *RedemptionDetails                    `json:"redemption,omitempty"`
	CreateData                           *CreateDataDetails                    `json:"create_data,omitempty"`
	UpdateData                           *UpdateDataDetails                    `json:"update_data,omitempty"`
	RemoveData                           *RemoveDataDetails                    `json:"remove_data,omitempty"`
	CreateDataCreationRequest            *CreateDataCreationRequest            `json:"create_data_creation_request,omitempty"`
	CancelDataCreationRequest            *CancelDataCreationRequest            `json:"cancel_data_creation_request,omitempty"`
	CreateDataUpdateRequest              *CreateDataUpdateRequest              `json:"create_data_update_request,omitempty"`
	CancelDataUpdateRequest              *CancelDataUpdateRequest              `json:"cancel_data_update_request,omitempty"`
	CreateDataRemoveRequest              *CreateDataRemoveRequest              `json:"create_data_remove_request,omitempty"`
	CancelDataRemoveRequest              *CancelDataRemoveRequest              `json:"cancel_data_remove_request,omitempty"`
	CreateDeferredPaymentCreationRequest *CreateDeferredPaymentCreationRequest `json:"create_deferred_payment_creation_request,omitempty"`
	CancelDeferredPaymentCreationRequest *CancelDeferredPaymentCreationRequest `json:"cancel_deferred_payment_creation_request,omitempty"`
	CreateCloseDeferredPaymentRequest    *CreateCloseDeferredPaymentRequest    `json:"create_close_deferred_payment_request,omitempty"`
	CancelCloseDeferredPaymentRequest    *CancelCloseDeferredPaymentRequest    `json:"cancel_close_deferred_payment_request,omitempty"`
	UpdateDataOwner                      *UpdateDataOwnerDetails               `json:"update_data_owner,omitempty"`
	CreateDataOwnerUpdateRequest         *CreateDataOwnerUpdateRequest         `json:"create_data_owner_update_request,omitempty"`
	CancelDataOwnerUpdateRequest         *CancelDataOwnerUpdateRequest         `json:"cancel_data_owner_update_request,omitempty"`
}

//Value - converts operation details into jsonb
func (r OperationDetails) Value() (driver.Value, error) {
	result, err := pgdb.JSONValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal operation details")
	}

	return result, nil
}

//Scan - converts jsonb into OperationDetails
func (r *OperationDetails) Scan(src interface{}) error {
	err := pgdb.JSONScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan operation details")
	}

	return nil
}

// CreateAccountDetails - stores details of create account operation
type CreateAccountDetails struct {
	AccountAddress string `json:"account_address"`
	AccountRole    uint64 `json:"account_role"`
}

//ManageBalanceDetails - details of ManageBalanceOp
type ManageBalanceDetails struct {
	DestinationAccount string                  `json:"destination_account"`
	Action             xdr.ManageBalanceAction `json:"action"`
	Asset              string                  `json:"asset"`
	BalanceAddress     string                  `json:"balance_address"`
}

// ManageAccountRuleDetails - details of ManageAccountRuleOp
type ManageAccountRuleDetails struct {
	Action        xdr.ManageAccountRuleAction `json:"action"`
	RuleID        uint64                      `json:"rule_id"`
	CreateDetails *UpdateAccountRuleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateAccountRuleDetails   `json:"update_details,omitempty"`
}

// UpdateAccountRuleDetails - details of new or updated rule
type UpdateAccountRuleDetails struct {
	Resource xdr.AccountRuleResource `json:"resource"`
	Action   xdr.AccountRuleAction   `json:"action"`
	IsForbid bool                    `json:"is_forbid"`
	Details  regources.Details       `json:"details"`
}

// ManageAccountRuleDetails - details of ManageAccountRuleOp
type ManageAccountSpecificRuleDetails struct {
	Action        xdr.ManageAccountSpecificRuleAction `json:"action"`
	RuleID        uint64                              `json:"rule_id"`
	CreateDetails *CreateAccountSpecificRuleDetails   `json:"create_details,omitempty"`
}

// UpdateAccountRuleDetails - details of new or updated rule
type CreateAccountSpecificRuleDetails struct {
	LedgerKey xdr.LedgerKey `json:"ledger_key"`
	Forbids   bool          `json:"forbids"`
	AccountID string        `json:"account_id"`
}

// ManageAccountRoleDetails - details of ManageAccountRoleOp
type ManageAccountRoleDetails struct {
	Action        xdr.ManageAccountRoleAction `json:"action"`
	RoleID        uint64                      `json:"role_id"`
	CreateDetails *UpdateAccountRoleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateAccountRoleDetails   `json:"update_details,omitempty"`
}

// UpdateAccountRoleDetails - details of new or updated role
type UpdateAccountRoleDetails struct {
	RuleIDs []uint64          `json:"rule_ids"`
	Details regources.Details `json:"details"`
}

// ManageSignerDetails - details op manage signer operation
type ManageSignerDetails struct {
	Action        xdr.ManageSignerAction `json:"action"`
	PublicKey     string                 `json:"public_key"`
	CreateDetails *UpdateSignerDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateSignerDetails   `json:"update_details,omitempty"`
}

// UpdateSignerDetails - details of new or updated signer
type UpdateSignerDetails struct {
	RoleID   uint64            `json:"role_id"`
	Weight   uint32            `json:"weight"`
	Identity uint32            `json:"identity"`
	Details  regources.Details `json:"details"`
}

// ManageSignerRuleDetails - details of ManageAccountRuleOp
type ManageSignerRuleDetails struct {
	Action        xdr.ManageSignerRuleAction `json:"action"`
	RuleID        uint64                     `json:"rule_id"`
	CreateDetails *CreateSignerRuleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateSignerRuleDetails   `json:"update_details,omitempty"`
}

// CreateSignerRuleDetails - details of new or updated rule
type CreateSignerRuleDetails struct {
	Resource   xdr.SignerRuleResource `json:"resource"`
	Action     xdr.SignerRuleAction   `json:"action"`
	IsForbid   bool                   `json:"is_forbid"`
	IsDefault  bool                   `json:"is_default"`
	IsReadOnly bool                   `json:"is_read_only"`
	Details    regources.Details      `json:"details"`
}

// UpdateSignerRuleDetails - details of new or updated rule
type UpdateSignerRuleDetails struct {
	Resource  xdr.SignerRuleResource `json:"resource"`
	Action    xdr.SignerRuleAction   `json:"action"`
	IsForbid  bool                   `json:"is_forbid"`
	IsDefault bool                   `json:"is_default"`
	Details   regources.Details      `json:"details"`
}

// ManageSignerRoleDetails - details of ManageAccountRoleOp
type ManageSignerRoleDetails struct {
	Action        xdr.ManageSignerRoleAction `json:"action"`
	RoleID        uint64                     `json:"role_id"`
	CreateDetails *CreateSignerRoleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateSignerRoleDetails   `json:"update_details,omitempty"`
}

// UpdateSignerRoleDetails - details of new or updated role
type CreateSignerRoleDetails struct {
	RuleIDs    []uint64          `json:"rule_ids"`
	IsReadOnly bool              `json:"is_read_only"`
	Details    regources.Details `json:"details"`
}

// UpdateSignerRoleDetails - details of new or updated role
type UpdateSignerRoleDetails struct {
	RuleIDs []uint64          `json:"rule_ids"`
	Details regources.Details `json:"details"`
}

//ManageKeyValueDetails - details of ManageKeyValueOp
type ManageKeyValueDetails struct {
	Action xdr.ManageKvAction            `json:"action"`
	Key    string                        `json:"key"`
	Value  *regources.KeyValueEntryValue `json:"value,omitempty"`
}

//SetFeeDetails - details of SetFeeOp
type SetFeeDetails struct {
	AccountAddress *string          `json:"account_address,omitempty"`
	AccountRole    *xdr.Uint64      `json:"account_role,omitempty"`
	AssetCode      string           `json:"asset_code"`
	FeeType        xdr.FeeType      `json:"fee_type"`
	FixedFee       regources.Amount `json:"fixed_fee"`
	IsDelete       bool             `json:"is_delete"`
	LowerBound     regources.Amount `json:"lower_bound"`
	PercentFee     regources.Amount `json:"percent_fee"`
	Subtype        int64            `json:"subtype"`
	UpperBound     regources.Amount `json:"upper_bound"`
	// FeeAsset deprecated
}

//CreateWithdrawRequestDetails - details of corresponding op
type CreateWithdrawRequestDetails struct {
	BalanceAddress string            `json:"balance_address"`
	Amount         regources.Amount  `json:"amount"`
	Fee            regources.Fee     `json:"fee"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateManageLimitsRequestDetails - details of corresponding op
type CreateManageLimitsRequestDetails struct {
	CreatorDetails regources.Details `json:"creator_details"`
	RequestID      int64             `json:"request_id"`
}

//ManageLimitsDetails - details of corresponding op
type ManageLimitsDetails struct {
	Action   xdr.ManageLimitsAction       `json:"action"`
	Creation *ManageLimitsCreationDetails `json:"creation_details,omitempty"`
	Removal  *ManageLimitsRemovalDetails  `json:"removal_details,omitempty"`
}

//ManageLimitsCreationDetails - details of corresponding op
type ManageLimitsCreationDetails struct {
	AccountRole     *xdr.Uint64      `json:"account_role,omitempty"`
	AccountAddress  string           `json:"account_address,omitempty"`
	StatsOpType     xdr.StatsOpType  `json:"stats_op_type"`
	AssetCode       string           `json:"asset_code"`
	IsConvertNeeded bool             `json:"is_convert_needed"`
	DailyOut        regources.Amount `json:"daily_out"`
	WeeklyOut       regources.Amount `json:"weekly_out"`
	MonthlyOut      regources.Amount `json:"monthly_out"`
	AnnualOut       regources.Amount `json:"annual_out"`
}

//ManageLimitsRemovalDetails - details of corresponding op
type ManageLimitsRemovalDetails struct {
	LimitsID int64 `json:"limits_id"`
}

//ManageAssetPairDetails - details of corresponding op
type ManageAssetPairDetails struct {
	BaseAsset               string              `json:"base_asset"`
	QuoteAsset              string              `json:"quote_asset"`
	PhysicalPrice           regources.Amount    `json:"physical_price"`
	PhysicalPriceCorrection regources.Amount    `json:"physical_price_correction"`
	MaxPriceStep            regources.Amount    `json:"max_price_step"`
	Policies                xdr.AssetPairPolicy `json:"policies"`
}

//ManageOfferDetails - details of corresponding op
type ManageOfferDetails struct {
	OfferID     int64            `json:"offer_id,omitempty"`
	OrderBookID int64            `json:"order_book_id"`
	BaseAsset   string           `json:"base_asset"`
	QuoteAsset  string           `json:"quote_asset"`
	Amount      regources.Amount `json:"base_amount"`
	Price       regources.Amount `json:"price"`
	IsBuy       bool             `json:"is_buy"`
	Fee         regources.Fee    `json:"fee"`
	IsDeleted   bool             `json:"is_deleted"`
}

//ReviewRequestDetails - details of corresponding op
type ReviewRequestDetails struct {
	Action          xdr.ReviewRequestOpAction         `json:"action"`
	Reason          string                            `json:"reason"`
	RequestHash     string                            `json:"request_hash"`
	RequestID       int64                             `json:"request_id"`
	RequestType     xdr.ReviewableRequestType         `json:"request_type"`
	RequestDetails  xdr.ReviewRequestOpRequestDetails `json:"request_details"`
	IsFulfilled     bool                              `json:"is_fulfilled"`
	AddedTasks      uint32                            `json:"added_tasks"`
	RemovedTasks    uint32                            `json:"removed_tasks"`
	ExternalDetails regources.Details                 `json:"external_details"`
}

//ManageAssetDetails - details of corresponding op
type ManageAssetDetails struct {
	AssetCode         string                `json:"asset_code,omitempty"`
	Type              uint64                `json:"type,omitempty"`
	RequestID         int64                 `json:"request_id"`
	Action            xdr.ManageAssetAction `json:"action"`
	Policies          *xdr.AssetPolicy      `json:"policies,omitempty"`
	CreatorDetails    regources.Details     `json:"creator_details,omitempty"`
	PreIssuanceSigner string                `json:"pre_issuance_signer,omitempty"`
	MaxIssuanceAmount regources.Amount      `json:"max_issuance_amount,omitempty"`
}

//CreatePreIssuanceRequestDetails - details of corresponding op
type CreatePreIssuanceRequestDetails struct {
	AssetCode      string            `json:"asset_code"`
	Amount         regources.Amount  `json:"amount"`
	CreatorDetails regources.Details `json:"creator_details"`
	RequestID      int64             `json:"request_id"`
	IsFulfilled    bool              `json:"is_fulfilled"`
}

//CreateIssuanceRequestDetails - details of corresponding op
type CreateIssuanceRequestDetails struct {
	Fee                    regources.Fee     `json:"fee"`
	Reference              string            `json:"reference"`
	Amount                 regources.Amount  `json:"amount"`
	Asset                  string            `json:"asset"`
	ReceiverAccountAddress string            `json:"receiver_account_address"`
	ReceiverBalanceAddress string            `json:"receiver_balance_address"`
	CreatorDetails         regources.Details `json:"creator_details"`
	AllTasks               *int64            `json:"all_tasks,omitempty"`
	RequestDetails         RequestDetails    `json:"request_details"`
}

//CreateSaleRequestDetails -details of corresponding op
type CreateSaleRequestDetails struct {
	RequestID         int64                  `json:"request_id"`
	BaseAsset         string                 `json:"base_asset"`
	DefaultQuoteAsset string                 `json:"default_quote_asset"`
	StartTime         time.Time              `json:"start_time"`
	EndTime           time.Time              `json:"end_time"`
	SoftCap           regources.Amount       `json:"soft_cap"`
	HardCap           regources.Amount       `json:"hard_cap"`
	QuoteAssets       []regources.AssetPrice `json:"quote_assets"`
	CreatorDetails    regources.Details      `json:"creator_details"`
}

//CheckSaleStateDetails - details of corresponding op
type CheckSaleStateDetails struct {
	SaleID int64 `json:"sale_id"`
	Effect xdr.CheckSaleStateEffect
}

//CreateAtomicSwapAskRequestDetails - details of corresponding op
type CreateAtomicSwapAskRequestDetails struct {
	Amount         regources.Amount       `json:"amount"`
	BaseBalance    string                 `json:"base_balance"`
	QuoteAssets    []regources.AssetPrice `json:"quote_assets"`
	Details        regources.Details      `json:"creator_details"`
	RequestDetails RequestDetails         `json:"request_details"`
}

//CreateAtomicSwapBidRequestDetails - details of corresponding op
type CreateAtomicSwapBidRequestDetails struct {
	AskID          int64             `json:"ask_id"`
	BaseAmount     regources.Amount  `json:"base_amount"`
	QuoteAsset     string            `json:"quote_asset"`
	RequestDetails RequestDetails    `json:"request_details"`
	Details        regources.Details `json:"creator_details"`
}

//CancelAtomicSwapAskDetails - details of corresponding op
type CancelAtomicSwapAskDetails struct {
	BidID int64 `json:"bid_id"`
}

//CreateAMLAlertRequestDetails - details of corresponding op
type CreateAMLAlertRequestDetails struct {
	Amount         regources.Amount  `json:"amount"`
	BalanceAddress string            `json:"balance_address"`
	CreatorDetails regources.Details `json:"creator_details"`
	RequestDetails RequestDetails    `json:"request_details"`
}

// PaymentDetails - stores details of payment operation
type PaymentDetails struct {
	AccountFrom             string           `json:"account_from"`
	AccountTo               string           `json:"account_to"`
	BalanceFrom             string           `json:"balance_from"`
	BalanceTo               string           `json:"balance_to"`
	Amount                  regources.Amount `json:"amount"`
	Asset                   string           `json:"asset"`
	SourceFee               regources.Fee    `json:"source_fee"`
	DestinationFee          regources.Fee    `json:"destination_fee"`
	SourcePayForDestination bool             `json:"source_pay_for_destination"`
	Subject                 string           `json:"subject"`
	Reference               string           `json:"reference"`
}

//PayoutDetails - details of corresponding op
type PayoutDetails struct {
	SourceAccountAddress string           `json:"source_account_address"`
	SourceBalanceAddress string           `json:"source_balance_address"`
	Asset                string           `json:"asset"`
	MaxPayoutAmount      regources.Amount `json:"max_payout_amount"`
	MinAssetHolderAmount regources.Amount `json:"min_asset_holder_amount"`
	MinPayoutAmount      regources.Amount `json:"min_payout_amount"`
	ExpectedFee          regources.Fee    `json:"expected_fee"`
	ActualFee            regources.Fee    `json:"actual_fee"`
	ActualPayoutAmount   regources.Amount `json:"actual_payout_amount"`
}

//RequestDetails - details of the request created or reviewed via op
type RequestDetails struct {
	RequestID   int64 `json:"request_id,omitempty"`
	IsFulfilled bool  `json:"is_fulfilled"`
}

//CreateChangeRoleRequestDetails - details of corresponding op
type CreateChangeRoleRequestDetails struct {
	DestinationAccount string            `json:"destination_account"`
	AccountRoleToSet   uint64            `json:"account_role_to_set"`
	CreatorDetails     regources.Details `json:"creator_details"`
	AllTasks           *uint32           `json:"all_tasks"`
	RequestDetails     RequestDetails    `json:"request_details"`
}

//ManageExternalSystemPoolDetails - details of corresponding op
type ManageExternalSystemPoolDetails struct {
	Action xdr.ManageExternalSystemAccountIdPoolEntryAction `json:"action"`
	Create *CreateExternalSystemPoolDetails                 `json:"create"`
	Remove *RemoveExternalSystemPoolDetails                 `json:"remove"`
}

//CreateExternalSystemPoolDetails - details of corresponding op
type CreateExternalSystemPoolDetails struct {
	Data               string `json:"data"`
	ExternalSystemType int32  `json:"external_system_type"`
	Parent             uint64 `json:"parent"`
	PoolId             uint64 `json:"pool_id"`
}

//RemoveExternalSystemPoolDetails - details of corresponding op
type RemoveExternalSystemPoolDetails struct {
	PoolId uint64 `json:"pool_id"`
}

//BindExternalSystemAccountDetails - details of corresponding op
type BindExternalSystemAccountDetails struct {
	ExternalSystemType int32 `json:"external_system_type"`
}

//ManageSaleDetails - details of corresponding op
type ManageSaleDetails struct {
	Action xdr.ManageSaleAction `json:"action"`
	SaleId uint64               `json:"sale_id"`
}

type LicenseDetails struct {
	PrevLicenseHash string    `json:"prev_license_hash"`
	AdminCount      uint64    `json:"admin_count"`
	DueDate         time.Time `json:"due_date"`
	LedgerHash      string    `json:"ledger_hash"`
	Signatures      []string  `json:"signatures"`
}

type StampDetails struct {
	LicenseHash string `json:"license_hash"`
	LedgerHash  string `json:"ledger_hash"`
}

type CreatePollRequestDetails struct {
	PollType                 xdr.PollType      `json:"poll_type"`
	PermissionType           uint32            `json:"permission_type"`
	NumberOfChoices          uint32            `json:"number_of_choices"`
	CreatorDetails           regources.Details `json:"creator_details"`
	StartTime                time.Time         `json:"start_time"`
	EndTime                  time.Time         `json:"end_time"`
	ResultProviderID         string            `json:"result_provider_id"`
	VoteConfirmationRequired bool              `json:"vote_confirmation_required"`
	PollData                 xdr.PollData      `json:"poll_data"`
	AllTasks                 *uint32           `json:"all_tasks,omitempty"`
	RequestDetails           RequestDetails    `json:"request_details"`
}

type CancelCreatePollRequestDetails struct {
	RequestID int64 `json:"request_id"`
}

type ManageCreatePollRequestDetails struct {
	Action        xdr.ManageCreatePollRequestAction `json:"action"`
	CreateDetails *CreatePollRequestDetails         `json:"create_details,omitempty"`
	CancelDetails *CancelCreatePollRequestDetails   `json:"cancel_details,omitempty"`
}

type ClosePollData struct {
	PollResult xdr.PollResult    `json:"poll_result"`
	Details    regources.Details `json:"details"`
}

type UpdatePollEndTimeData struct {
	EndTime time.Time `json:"new_end_time"`
}

type ManagePollDetails struct {
	PollID            int64                  `json:"poll_id"`
	Action            xdr.ManagePollAction   `json:"action"`
	ClosePoll         *ClosePollData         `json:"close_poll,omitempty"`
	UpdatePollEndTime *UpdatePollEndTimeData `json:"update_poll_end_time,omitempty"`
}

type ManageVoteDetails struct {
	PollID   int64                `json:"poll_id"`
	Action   xdr.ManageVoteAction `json:"action"`
	VoteData *xdr.VoteData        `json:"vote_data,omitmepty"`
}

type RemoveAssetPairDetails struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type InitiateKYCRecoveryDetails struct {
	Account string `json:"account"`
	Signer  string `json:"signer"`
}

type CreateKYCRecoveryRequestDetails struct {
	TargetAccount  string                       `json:"target_account"`
	SignersData    []regources.UpdateSignerData `json:"signers_data"`
	CreatorDetails regources.Details            `json:"creator_details"`
	AllTasks       *uint32                      `json:"all_tasks"`
	RequestDetails RequestDetails               `json:"request_details"`
}

type CreateManageOfferRequestDetails struct {
	ManageOfferDetails ManageOfferRequest `json:"manage_offer_details"`
	AllTasks           *uint32            `json:"all_tasks"`
	RequestDetails     RequestDetails     `json:"request_details"`
}

type PaymentRequestDetails struct {
	CreatorDetails          regources.Details `json:"creator_details"`
	AccountFrom             string            `json:"account_from"`
	BalanceFrom             string            `json:"balance_from"`
	Amount                  regources.Amount  `json:"amount"`
	SourceFee               regources.Fee     `json:"source_fee"`
	DestinationFee          regources.Fee     `json:"destination_fee"`
	SourcePayForDestination bool              `json:"source_pay_for_destination"`
	Subject                 string            `json:"subject"`
	Reference               string            `json:"reference"`
}

type CreatePaymentRequestDetails struct {
	PaymentDetails PaymentRequestDetails `json:"payment_details"`
	AllTasks       *uint32               `json:"all_tasks"`
	RequestDetails RequestDetails        `json:"request_details"`
}

type RemoveAssetDetails struct {
	Code string `json:"code"`
}

type OpenSwapDetails struct {
	AccountFrom             string            `json:"account_from"`
	AccountTo               string            `json:"account_to"`
	BalanceFrom             string            `json:"balance_from"`
	BalanceTo               string            `json:"balance_to"`
	Amount                  regources.Amount  `json:"amount"`
	Asset                   string            `json:"asset"`
	SourceFee               regources.Fee     `json:"source_fee"`
	DestinationFee          regources.Fee     `json:"destination_fee"`
	SourcePayForDestination bool              `json:"source_pay_for_destination"`
	SecretHash              string            `json:"secret_hash"`
	LockTime                time.Time         `json:"lock_time"`
	Details                 regources.Details `db:"details"`
}

type CloseSwapDetails struct {
	ID     int64   `json:"id"`
	Secret *string `json:"secret"`
}

type RedemptionDetails struct {
	BalanceFrom    string            `json:"balance_from"`
	AccountTo      string            `json:"account_to"`
	Amount         regources.Amount  `json:"amount"`
	Asset          string            `json:"asset"`
	Details        regources.Details `json:"details"`
	RequestDetails RequestDetails    `json:"request_details"`
}

type CreateDataDetails struct {
	ID    uint64            `json:"id"`
	Type  uint64            `json:"type"`
	Value regources.Details `json:"value"`
	Owner string            `json:"owner"`
}

type UpdateDataDetails struct {
	ID    uint64            `json:"id"`
	Value regources.Details `json:"value"`
}

type UpdateDataOwnerDetails struct {
	ID       uint64        `json:"id"`
	NewOwner xdr.AccountId `json:"newOwner"`
}

type RemoveDataDetails struct {
	ID uint64 `json:"id"`
}

type CreateDataCreationRequest struct {
	ID             uint64            `json:"id"`
	Type           uint64            `json:"type"`
	Value          regources.Details `json:"value"`
	Owner          string            `json:"owner"`
	CreatorDetails regources.Details `json:"creator_details"`
	RequestID      uint64            `json:"request_id"`
}

type CreateDataUpdateRequest struct {
	ID             uint64            `json:"id"`
	Value          regources.Details `json:"value"`
	CreatorDetails regources.Details `json:"creator_details"`
	RequestID      uint64            `json:"request_id"`
}

type CreateDataOwnerUpdateRequest struct {
	ID             uint64            `json:"id"`
	NewOwner       xdr.AccountId     `json:"new_owner"`
	CreatorDetails regources.Details `json:"creator_details"`
	RequestID      uint64            `json:"request_id"`
}

type CreateDataRemoveRequest struct {
	ID             uint64            `json:"id"`
	CreatorDetails regources.Details `json:"creator_details"`
	RequestID      uint64            `json:"request_id"`
}

type CancelDataCreationRequest struct {
	RequestID uint64 `json:"request_id"`
}

type CancelDataUpdateRequest struct {
	RequestID uint64 `json:"request_id"`
}

type CancelDataOwnerUpdateRequest struct {
	RequestID uint64 `json:"request_id"`
}

type CancelDataRemoveRequest struct {
	RequestID uint64 `json:"request_id"`
}

type CreateDeferredPaymentCreationRequest struct {
	RequestID          uint64            `json:"request_id"`
	SourceBalance      string            `json:"source_balance"`
	DestinationAccount string            `json:"destination_account"`
	Amount             regources.Amount  `json:"amount"`
	Details            regources.Details `json:"details"`
	AllTasks           *uint32           `json:"all_tasks"`
}

type CreateCloseDeferredPaymentRequest struct {
	RequestID          uint64            `json:"request_id"`
	DestinationBalance string            `json:"destination_balance"`
	DestinationAccount string            `json:"destination_account"`
	DeferredPaymentID  uint64            `json:"deferred_payment_id"`
	Amount             regources.Amount  `json:"amount"`
	Details            regources.Details `json:"details"`
	AllTasks           *uint32           `json:"all_tasks"`
}

type CancelDeferredPaymentCreationRequest struct {
	RequestID uint64 `json:"request_id"`
}

type CancelCloseDeferredPaymentRequest struct {
	RequestID uint64 `json:"request_id"`
}
