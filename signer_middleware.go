package horizon

import (
	"net/http"

	"github.com/SafeRE-IT/horizon/render/problem"
	"github.com/zenazn/goji/web"
	"gitlab.com/tokend/go/signcontrol"
)

type SignatureValidator struct {
	SkipCheck bool
}

func (v *SignatureValidator) GetSigner(r *http.Request) (*string, error) {
	if v.SkipCheck {
		return nil, nil
	}

	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		return nil, err
	}

	if signer == "" {
		return nil, nil
	}

	return &signer, nil
}

// Middleware checks only if request signature is valid and sets signer to request
// context if true
func (v *SignatureValidator) Middleware(c *web.C, next http.Handler) http.Handler {
	return v.MiddlewareV2(next)
}

func (v *SignatureValidator) MiddlewareV2(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, err := v.GetSigner(r)
		if err != nil {
			switch err {
			case signcontrol.ErrNotSigned:
				// passing not signed requests through w/o setting any headers
				next.ServeHTTP(w, r)
			default:
				p := &problem.BadRequest
				p.Extras = map[string]interface{}{
					"invalid_field": "signature",
					"reason":        err.Error(),
				}
				problem.Render(r.Context(), w, p)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
