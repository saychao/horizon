package requests

import (
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/go-chi/chi"
)

type GetDeferredPayment struct {
	*base
	DeferredPaymentID int64
}

func NewGetDeferredPayment(r *http.Request) (*GetDeferredPayment, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeDeferredPaymentAll,
	})

	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert string to int64")
	}

	req := GetDeferredPayment{
		base:   b,
		DeferredPaymentID: id,
	}

	return &req, nil
}
