// Package sqx contains utilities and extensions for the squirrel package which
// is used by horizon to generate sql statements.
package sqx

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/tokend/go/xdr"
)

// StringArray returns a sq.Expr suitable for inclusion in an insert that represents
// the Postgres-compatible array insert.
func StringArray(str []string) interface{} {
	return sq.Expr(
		"?::character varying[]",
		fmt.Sprintf("{%s}", strings.Join(str, ",")),
	)
}

func InForReviewableRequestTypes(columnName string, values ...xdr.ReviewableRequestType) (string, []interface{}) {
	rawValues := make([]interface{}, len(values))
	for i := range values {
		rawValues[i] = int32(values[i])
	}

	return In(columnName, rawValues...)
}

// Returns statement and params of it for SQL IN.
func InForString(columnName string, values ...string) (string, []interface{}) {
	rawValues := make([]interface{}, len(values))
	for i := range values {
		rawValues[i] = values[i]
	}

	return In(columnName, rawValues...)
}

// Returns statement and params of it for SQL IN.
func In(columnName string, values ...interface{}) (string, []interface{}) {
	params := make([]string, len(values))
	for i := range values {
		params[i] = "?"
	}

	return fmt.Sprintf("%s IN (%s)", columnName, strings.Join(params, ",")), values
}
