package txsub

import (
	"time"

	"github.com/SafeRE-IT/horizon/corer"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/net/context"
)

// NewDefaultSubmitter returns a new, simple Submitter implementation
// that submits directly to the core at `url` using the http client
// `h`.
func NewDefaultSubmitter(connector corer.Connector) Submitter {
	return &submitter{
		connector: connector,
	}
}

// submitter is the default implementation for the Submitter interface.  It
// submits directly to the configured core instance using the
// configured http client.
type submitter struct {
	connector corer.Connector
}

// Submit sends the provided envelope to core and parses the response into
// a SubmissionResult
func (sub *submitter) Submit(ctx context.Context, env *EnvelopeInfo) (duration time.Duration, err error) {
	start := time.Now()
	defer func() { duration = time.Since(start) }()

	coreResponse, err := sub.connector.SubmitTx(env.RawBlob)
	if err != nil {
		err = errors.Wrap(err, "Failed to submit tx to core")
		return
	}

	// interpet response
	if coreResponse.Exception != "" {
		err = errors.From(errors.New("Core exception"), logan.F{
			"exception": coreResponse.Exception,
		})
		return
	}

	switch coreResponse.Status {
	case "ERROR":
		err = NewRejectedTxError(coreResponse.Error)
	case "PENDING", "DUPLICATE":
	//noop.  A nil Err indicates success
	default:
		err = errors.From(errors.New("Unrecognized core status response"), logan.F{
			"status": coreResponse.Status,
		})
	}

	return
}
