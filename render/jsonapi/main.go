package jsonapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/saychao/horizon/log"
	"github.com/saychao/horizon/render/problem"
)

type ErrorObject struct {
	// ID is a unique identifier for this particular occurrence of a problem.
	ID string `json:"id,omitempty"`

	// Title is a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence
	// of the problem, except for purposes of localization.
	Title string `json:"title,omitempty"`

	// Detail is a human-readable explanation specific to this occurrence of the problem. Like title, this field’s
	// value can be localized.
	Detail string `json:"detail,omitempty"`

	// Status is the HTTP status code applicable to this problem, expressed as a string value.
	Status string `json:"status,omitempty"`

	// Code is an application-specific error code, expressed as a string value.
	Code string `json:"code,omitempty"`

	// Meta is an object containing non-standard meta-information about the error.
	Meta map[string]interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Errors []*ErrorObject `json:"errors"`
}

var (
	errToJsonApiMap = map[error]ErrorObject{}
)

// RegisterError records an error -> ErrorObject mapping, allowing the app to register
// specific errors that may occur in other packages to be rendered as a specific
// ErrorObject instance.
func RegisterError(err error, errObj ErrorObject) {
	errToJsonApiMap[err] = errObj
}

func RenderErr(ctx context.Context, w http.ResponseWriter, err interface{}) {
	switch err := err.(type) {
	case problem.P:
		render(ctx, w, FromProblem(err))
	case *problem.P:
		render(ctx, w, FromProblem(*err))
	case problem.HasProblem:
		render(ctx, w, FromProblem(err.Problem()))
	case error:
		render(ctx, w, FromError(err))
	default:
		panic(fmt.Sprintf("Invalid error: %v+", err))
	}
}

func render(ctx context.Context, w http.ResponseWriter, errObjects ...*ErrorObject) {
	w.Header().Set("Content-Type", "application/vnd.api+json")

	js, err := json.MarshalIndent(
		ErrorResponse{
			Errors: errObjects,
		},
		"",
		"  ",
	)

	if err != nil {
		err := errors.Wrap(err, 1)
		log.Ctx(ctx).WithStack(err).Error(err)
		http.Error(w, "error rendering problem", http.StatusInternalServerError)
		return
	}

	status, err := strconv.Atoi(errObjects[0].Status)
	if err != nil {
		panic(fmt.Sprintf("Invalid status: %s+", errObjects[0].Status))
	}

	w.WriteHeader(status)
	w.Write(js)
}
