package changes

import (
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

type accountStatusStorage interface {
	SetKYCRecoveryStatus(address string, status int) error
}

type signerHandler struct {
	accountStatusStorage accountStatusStorage
}

func newSignerHandler(storage accountStatusStorage) *signerHandler {
	return &signerHandler{
		accountStatusStorage: storage,
	}
}

//We don't care about other causes
func (p *signerHandler) Created(lc ledgerChange) error {
	op := lc.Operation
	accID := lc.LedgerChange.MustCreated().Data.MustSigner().AccountId
	switch op.Body.Type {
	case xdr.OperationTypeReviewRequest:
		reviewRequestOp := op.Body.MustReviewRequestOp()
		switch reviewRequestOp.RequestDetails.RequestType {
		case xdr.ReviewableRequestTypeKycRecovery:
			return p.accountStatusStorage.SetKYCRecoveryStatus(accID.Address(),
				int(regources.KYCRecoveryStatusNone))
		}
	case xdr.OperationTypeInitiateKycRecovery:
		initKycRecovery := op.Body.MustInitiateKycRecoveryOp()
		return p.accountStatusStorage.SetKYCRecoveryStatus(initKycRecovery.Account.Address(),
			int(regources.KYCRecoveryStatusInitiated))
	}

	return nil
}
