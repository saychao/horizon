package codes

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/tokend/go/xdr"
)

func TestCodeProviders(t *testing.T) {
	Convey("Code providers", t, func() {
		for _, opType := range xdr.OperationTypeAll {
			_, ok := codeProviders[opType]
			if !ok {
				t.Errorf("Failed to get code provider for %s", opType.ShortString())
			}
		}
	})
}
