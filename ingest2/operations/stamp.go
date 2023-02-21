package operations

import (
	"encoding/hex"

	"gitlab.com/tokend/go/xdr"
	"github.com/SafeRE-IT/horizon/db2/history2"
)

type stampOpHandler struct {
	effectsProvider
}

// Details returns details about manage asset operation
func (h *stampOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	stampResult := opRes.StampResult.MustSuccess()
	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeStamp,
		Stamp: &history2.StampDetails{
			LedgerHash:  hex.EncodeToString([]byte(stampResult.LedgerHash[:])),
			LicenseHash: hex.EncodeToString([]byte(stampResult.LicenseHash[:])),
		},
	}

	return opDetails, nil
}
