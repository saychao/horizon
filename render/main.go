package render

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/log"
	"github.com/SafeRE-IT/horizon/render/jsonapi"
	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/munnerz/goautoneg"
	"golang.org/x/net/context"
)

// Negotiate inspects the Accept header of the provided request and determines
// what the most appropriate response type should be.  Defaults to HAL.
func Negotiate(ctx context.Context, r *http.Request) string {
	alternatives := []string{MimeHal, MimeJSON, MimeEventStream, MimeRaw, MimeJSONAPI}
	accept := r.Header.Get("Accept")

	if accept == "" {
		return MimeHal
	}

	result := goautoneg.Negotiate(r.Header.Get("Accept"), alternatives)

	log.Ctx(ctx).WithFields(log.F{
		"content_type": result,
		"accept":       accept,
	}).Debug("Negotiated content type")

	return result
}

func RegisterError(err error, p problem.P) {
	problem.RegisterError(err, p)
	jsonapi.RegisterError(err, *jsonapi.FromProblem(p))
}
