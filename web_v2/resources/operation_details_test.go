package resources

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/tokend/go/xdr"
)

func TestOperationDetailsProviders(t *testing.T) {
	Convey("Operation details providers", t, func() {
		for _, opType := range xdr.OperationTypeAll {
			_, ok := operationDetailsProviders[opType]
			if !ok {
				t.Errorf("Failed to get operation details provider for %s", opType.ShortString())
			}
		}
	})
}
