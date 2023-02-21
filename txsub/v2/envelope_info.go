package txsub

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/build"
	"gitlab.com/tokend/go/xdr"
)

// EnvelopeInfo provides information needed for tx envelope submission
type EnvelopeInfo struct {
	// Hash of the transaction
	ContentHash string
	// Source account of the tx
	SourceAddress string
	// base64 encoded envelope
	RawBlob string
}

func ExtractEnvelopeInfo(env string, passphrase string) (*EnvelopeInfo, error) {
	var tx xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(env, &tx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal tx")
	}

	txb := build.TransactionBuilder{TX: &tx.Tx}
	txb.Mutate(build.Network{passphrase})

	var result EnvelopeInfo
	result.ContentHash, err = txb.HashHex()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get hash of tx")
	}

	result.SourceAddress = tx.Tx.SourceAccount.Address()
	result.RawBlob = env

	return &result, nil
}

func (envelopeInfo EnvelopeInfo) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"content_hash":   envelopeInfo.ContentHash,
		"source_address": envelopeInfo.SourceAddress,
		"raw_blob":       envelopeInfo.RawBlob,
	}
}
