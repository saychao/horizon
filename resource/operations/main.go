package operations

import (
	"github.com/saychao/horizon/db2/history"
	"github.com/saychao/horizon/render/hal"
	"gitlab.com/tokend/go/xdr"
	"golang.org/x/net/context"
)

// New creates a new operation resource, finding the appropriate type to use
// based upon the row's type.
func New(
	ctx context.Context, row history.Operation, participants []*history.Participant, public bool,
) (result hal.Pageable, err error) {

	base := Base{}
	err = base.Populate(ctx, row, participants, public)
	if err != nil {
		return
	}
	switch row.Type {
	case xdr.OperationTypeCreateAccount:
		e := CreateAccount{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.Funder = ""
			e.Account = ""
			e.Referrer = nil
		}
		result = e
	case xdr.OperationTypeSetFees:
		e := SetFees{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			if e.Fee != nil {
				e.Fee.AccountID = ""
			}
		}
		result = e
	case xdr.OperationTypeCreateWithdrawalRequest:
		e := CreateWithdrawalRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.ExternalDetails = nil
		}
		result = e
	case xdr.OperationTypeManageLimits:
		e := ManageLimits{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageInvoiceRequest:
		e := ManageInvoiceRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageOffer:
		e := ManagerOffer{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageAssetPair:
		e := ManageAssetPair{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeCreateIssuanceRequest:
		e := CreateIssuanceRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.ExternalDetails = nil
		}
		result = e
	case xdr.OperationTypeCheckSaleState:
		e := CheckSaleState{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeCreateAmlAlert:
		e := CreateAmlAlert{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.BalanceID = ""
		}
		result = e
	case xdr.OperationTypeCreateChangeRoleRequest:
		e := CreateChangeRoleRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		if public {
			e.KYCData = nil
		}
		result = e
	case xdr.OperationTypeReviewRequest:
		e := ReviewRequest{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypePayment:
		e := PaymentV2{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageSale:
		e := ManageSale{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	case xdr.OperationTypeManageAsset:
		e := ManageAsset{Base: base}
		err = row.UnmarshalDetails(&e)
		result = e
	default:
		result = base
	}

	return
}

// CreateAccount is the json resource representing a single operation whose type
// is CreateAccount.
type CreateAccount struct {
	Base
	Funder      string  `json:"funder,omitempty"`
	Account     string  `json:"account,omitempty"`
	AccountType int32   `json:"account_type"`
	Referrer    *string `json:"referrer,omitempty"`
}

type BasePayment struct {
	From                  string `json:"from,omitempty"`
	To                    string `json:"to,omitempty"`
	FromBalance           string `json:"from_balance,omitempty"`
	ToBalance             string `json:"to_balance,omitempty"`
	Amount                string `json:"amount"`
	UserDetails           string `json:"user_details,omitempty"`
	Asset                 string `json:"asset"`
	SourcePaymentFee      string `json:"source_payment_fee"`
	DestinationPaymentFee string `json:"destination_payment_fee"`
	SourceFixedFee        string `json:"source_fixed_fee"`
	DestinationFixedFee   string `json:"destination_fixed_fee"`
	SourcePaysForDest     bool   `json:"source_pays_for_dest"`
}

// Payment is the json resource representing a single operation whose type is
// Payment.
type Payment struct {
	Base
	BasePayment
	Subject   string `json:"subject,omitempty"`
	Reference string `json:"reference,omitempty"`
	Asset     string `json:"qasset"`
}

// SetOptions is the json resource representing a single operation whose type is
// SetOptions.
type SetOptions struct {
	Base
	HomeDomain    string `json:"home_domain,omitempty"`
	InflationDest string `json:"inflation_dest,omitempty"`

	MasterKeyWeight *int   `json:"master_key_weight,omitempty"`
	SignerKey       string `json:"signer_key,omitempty"`
	SignerWeight    *int   `json:"signer_weight,omitempty"`
	SignerType      *int   `json:"signer_type,omitempty"`
	SignerIdentity  *int   `json:"signer_identity,omitempty"`

	SetFlags    []int    `json:"set_flags,omitempty"`
	SetFlagsS   []string `json:"set_flags_s,omitempty"`
	ClearFlags  []int    `json:"clear_flags,omitempty"`
	ClearFlagsS []string `json:"clear_flags_s,omitempty"`

	LowThreshold  *int `json:"low_threshold,omitempty"`
	MedThreshold  *int `json:"med_threshold,omitempty"`
	HighThreshold *int `json:"high_threshold,omitempty"`
}

//SetFees is the json resource representing a single operation whose type
// is SetFees.

type Fee struct {
	AssetCode   string `json:"asset_code"`
	FixedFee    string `json:"fixed_fee"`
	PercentFee  string `json:"percent_fee"`
	FeeType     int64  `json:"fee_type"`
	AccountID   string `json:"account_id,omitempty"`
	AccountType int64  `json:"account_type"`
	Subtype     int64  `json:"subtype"`
	LowerBound  int64  `json:"lower_bound"`
	UpperBound  int64  `json:"upper_bound"`
}

type SetFees struct {
	Base
	Fee *Fee `json:"fee"`
}

type ManagerOffer struct {
	Base
	IsBuy       bool   `json:"is_buy"`
	BaseAsset   string `json:"base_asset"`
	Amount      string `json:"amount"`
	Price       string `json:"price"`
	Fee         string `json:"fee"`
	OfferId     int64  `json:"offer_id"`
	OrderBookID int64  `json:"order_book_id"`
	IsDeleted   bool   `json:"is_deleted"`
}
