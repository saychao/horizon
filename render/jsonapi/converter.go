package jsonapi

import (
	"strconv"

	"github.com/SafeRE-IT/horizon/render/problem"
)

func FromProblem(p problem.P) *ErrorObject {
	errorObject := ErrorObject{
		Title:  p.Title,
		Detail: p.Detail,
		Status: strconv.Itoa(p.Status),
		Code:   p.Type,
	}

	if p.Extras != nil {
		errorObject.Meta = map[string]interface{}{
			"extras": p.Extras,
		}
	}

	return &errorObject
}

func FromError(err error) *ErrorObject {
	errObj, ok := errToJsonApiMap[err]

	// If this error is not a registered error
	// log it and replace it with a 500 error
	if !ok {
		return FromProblem(problem.ServerError)
	}

	return &errObj
}
