package operations

import (
	"testing"

	"gitlab.com/tokend/go/xdr"
)

func TestAllReviewableRequestsHandled(t *testing.T) {
	reviewRequestOpHandlerInst := newReviewRequestOpHandler(effectsProvider{&MockIDProvider{}, &mockBalanceProvider{}}, nil)
	for _, requestT := range xdr.ReviewableRequestTypeAll {
		if requestT == xdr.ReviewableRequestTypeNone || requestT == xdr.ReviewableRequestTypeAny {
			continue
		}
		_, ok := reviewRequestOpHandlerInst.allRequestHandlers[requestT]
		if !ok {
			t.Fatalf("All reivable requests must be handled. Request type: %s is not handled", requestT)
		}
	}
}
