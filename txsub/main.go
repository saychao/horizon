package txsub

import (
	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/txsub"
	"gitlab.com/tokend/go/build"
	"gitlab.com/tokend/go/xdr"
)

func ExtractEnvelopeInfo(env string, passphrase string) (*txsub.EnvelopeInfo, error) {
	var tx xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(env, &tx)
	if err != nil {
		return nil, logan.Wrap(err, "Failed to unmarshal tx")
	}

	txb := build.TransactionBuilder{TX: &tx.Tx}
	txb.Mutate(build.Network{passphrase})

	var result txsub.EnvelopeInfo
	result.ContentHash, err = txb.HashHex()
	if err != nil {
		return nil, logan.Wrap(err, "Failed to get hash of tx")
	}

	result.SourceAddress = tx.Tx.SourceAccount.Address()
	result.RawBlob = env

	return &result, nil
}
