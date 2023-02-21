package operations

import (
	"encoding/hex"

	"github.com/SafeRE-IT/horizon/db2/history2"
	"github.com/SafeRE-IT/horizon/ingest2/internal"
	"gitlab.com/tokend/go/xdr"
)

type licenseOpHandler struct {
	effectsProvider
}

// Details returns details about manage asset operation
func (h *licenseOpHandler) Details(op rawOperation, opRes xdr.OperationResultTr,
) (history2.OperationDetails, error) {
	licenseOp := op.Body.MustLicenseOp()
	signatures := make([]string, 0, len(licenseOp.Signatures))
	for _, v := range licenseOp.Signatures {
		signatures = append(signatures, hex.EncodeToString(v.Signature))
	}
	opDetails := history2.OperationDetails{
		Type: xdr.OperationTypeLicense,
		License: &history2.LicenseDetails{
			AdminCount:      uint64(licenseOp.AdminCount),
			DueDate:         internal.TimeFromXdr(licenseOp.DueDate),
			LedgerHash:      hex.EncodeToString([]byte(licenseOp.LedgerHash[:])),
			PrevLicenseHash: hex.EncodeToString([]byte(licenseOp.PrevLicenseHash[:])),
			Signatures:      signatures,
		},
	}

	return opDetails, nil
}
