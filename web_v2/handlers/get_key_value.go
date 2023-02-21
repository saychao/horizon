package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/db2/core2"
	"github.com/SafeRE-IT/horizon/web_v2/ctx"
	"github.com/SafeRE-IT/horizon/web_v2/requests"
	"github.com/SafeRE-IT/horizon/web_v2/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	regources "gitlab.com/tokend/regources/generated"
)

func GetKeyValue(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetKeyValue(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	keyValueQ := core2.NewKeyValueQ(ctx.CoreRepo(r))

	record, err := keyValueQ.ByKey(request.Key)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get key value")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if record == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resource := resources.NewKeyValue(*record)
	response := &regources.KeyValueEntryResponse{
		Data: resource,
	}

	ape.Render(w, response)
}
