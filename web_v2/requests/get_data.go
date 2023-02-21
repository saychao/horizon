package requests

import (
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/go-chi/chi"
)

type GetData struct {
	*base
	DataID int64
}

func NewGetData(r *http.Request) (*GetData, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeDataAll,
	})

	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert string to int64")
	}

	req := GetData{
		base:   b,
		DataID: id,
	}

	return &req, nil
}
