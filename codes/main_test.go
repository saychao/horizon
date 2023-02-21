package codes

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/tokend/go/xdr"
)

func checkCode(t *testing.T, code shortStr) {
	opCode := opCodeToString(code)
	So(opCode, ShouldNotBeBlank)
	message := getMessage(opCode)
	if message == "" {
		t.Errorf("Expected message not to be blank for %s", opCode)
	}
}

func TestCodes(t *testing.T) {
	Convey("TransactionResultCode", t, func() {
		for _, code := range xdr.TransactionResultCodeAll {
			message := getMessage(code.ShortString())
			if message == "" || message == code.ShortString() {
				t.Errorf("Expected message not to be blank for %s", code.ShortString())
			}
		}
	})
	Convey("OperationResultCode", t, func() {
		for _, code := range xdr.OperationResultCodeAll {
			message := getMessage(code.ShortString())
			if message == "" || message == code.ShortString() {
				t.Errorf("Expected message not to be blank for %s", code.ShortString())
			}
		}
	})
	Convey("CreateAccountResultCode", t, func() {
		for _, code := range xdr.CreateAccountResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageSignerResultCode", t, func() {
		for _, code := range xdr.ManageSignerResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("SetFeesResultCode", t, func() {
		for _, code := range xdr.SetFeesResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageSignerRoleResultCode", t, func() {
		for _, code := range xdr.ManageSignerRoleResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageSignerRuleResultCode", t, func() {
		for _, code := range xdr.ManageSignerRuleResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("CreateWithdrawalRequestResultCode", t, func() {
		for _, code := range xdr.CreateWithdrawalRequestResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("CreateIssuanceRequestResultCode", t, func() {
		for _, code := range xdr.CreateIssuanceRequestResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageBalanceResultCode", t, func() {
		for _, code := range xdr.ManageBalanceResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("ManageAssetResultCode", t, func() {
		for _, code := range xdr.ManageAssetResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageLimits", t, func() {
		for _, code := range xdr.ManageLimitsResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageAssetPair", t, func() {
		for _, code := range xdr.ManageAssetPairResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageOffer", t, func() {
		for _, code := range xdr.ManageOfferResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("ManageInvoiceRequest", t, func() {
		for _, code := range xdr.ManageInvoiceRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Withdrawal", t, func() {
		for _, code := range xdr.CreateWithdrawalRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create sale", t, func() {
		for _, code := range xdr.CreateSaleCreationRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Review request", t, func() {
		for _, code := range xdr.ReviewRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Manage external system account id pool entry", t, func() {
		for _, code := range xdr.ManageExternalSystemAccountIdPoolEntryResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Bind external system account id", t, func() {
		for _, code := range xdr.BindExternalSystemAccountIdResultCodeAll {
			checkCode(t, code)
		}
	})
	Convey("Create KYC request", t, func() {
		for _, code := range xdr.CreateChangeRoleRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Payment V2", t, func() {
		for _, code := range xdr.PaymentResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Manage sale", t, func() {
		for _, code := range xdr.ManageSaleResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Manage key value", t, func() {
		for _, code := range xdr.ManageKeyValueResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create manage limits request", t, func() {
		for _, code := range xdr.ManageLimitsResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Manage contract request", t, func() {
		for _, code := range xdr.ManageContractRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Manage contract", t, func() {
		for _, code := range xdr.ManageContractResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create Atomic Swap Bid Creation Request", t, func() {
		for _, code := range xdr.CreateAtomicSwapBidRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Cancel atomic swap bid", t, func() {
		for _, code := range xdr.CancelAtomicSwapAskResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create Atomic Swap Request", t, func() {
		for _, code := range xdr.CreateAtomicSwapAskRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Cancel CreateSale Request", t, func() {
		for _, code := range xdr.CancelSaleCreationRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Remove asset pair", t, func() {
		for _, code := range xdr.RemoveAssetPairResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create Manage Offer Request", t, func() {
		for _, code := range xdr.CreateManageOfferRequestResultCodeAll {
			checkCode(t, code)
		}
	})

	Convey("Create Payment Request", t, func() {
		for _, code := range xdr.CreatePaymentRequestResultCodeAll {
			checkCode(t, code)
		}
	})
}
