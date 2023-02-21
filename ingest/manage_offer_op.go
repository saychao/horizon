package ingest

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (is *Session) updateOfferState(offerID, state uint64) error {

	err := is.Ingestion.UpdateOfferState(offerID, state)
	if err != nil {
		return errors.Wrap(err, "failed to update offer state")
	}
	return nil
}
