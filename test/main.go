// Package test contains simple test helpers that should not
// have any dependencies on horizon's packages.  think constants,
// custom matchers, generic helpers etc.
package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	hlog "github.com/saychao/horizon/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

// StaticMockServer is a test helper that records it's last request
type StaticMockServer struct {
	*httptest.Server
	LastRequest *http.Request
}

// T provides a common set of functionality for each test in horizon
type T struct {
	T          *testing.T
	Assert     *assert.Assertions
	Require    *require.Assertions
	Ctx        context.Context
	HorizonDB  *sqlx.DB
	CoreDB     *sqlx.DB
	Logger     *hlog.Entry
	LogMetrics *hlog.Metrics
	LogBuffer  *bytes.Buffer
}

// Context provides a context suitable for testing in tests that do not create
// a full App instance (in which case your tests should be using the app's
// context).  This context has a logger bound to it suitable for testing.
func Context() context.Context {
	return hlog.Set(context.Background(), testLogger)
}

// ContextWithLogBuffer returns a context and a buffer into which the new, bound
// logger will write into.  This method allows you to inspect what data was
// logged more easily in your tests.
func ContextWithLogBuffer() (context.Context, *bytes.Buffer) {
	output := new(bytes.Buffer)
	l, _ := hlog.New()
	l.Entry.Out(output)
	l.Formatter(&logrus.TextFormatter{DisableColors: true})
	l.Level(hlog.DebugLevel)

	ctx := hlog.Set(context.Background(), l)
	return ctx, output

}

// OverrideLogger sets the default logger used by horizon to `l`.  This is used
// by the testing system so that we can collect output from logs during test
// runs.  Panics if the logger is already overridden.
func OverrideLogger(l *hlog.Entry) {
	if oldDefault != nil {
		panic("logger already overridden")
	}

	oldDefault = hlog.DefaultLogger
	hlog.DefaultLogger = l
}

// RestoreLogger restores the default horizon logger after it is overridden
// using a call to `OverrideLogger`.  Panics if the default logger is not
// presently overridden.
func RestoreLogger() {
	if oldDefault == nil {
		panic("logger not overridden, cannot restore")
	}

	hlog.DefaultLogger = oldDefault
	oldDefault = nil
}

// Start initializes a new test helper object and conceptually "starts" a new
// test
func Start(t *testing.T) *T {
	result := &T{}

	result.T = t
	result.LogBuffer = new(bytes.Buffer)
	result.Logger, result.LogMetrics = hlog.New()
	result.Logger.Out(result.LogBuffer)
	result.Logger.Formatter(&logrus.TextFormatter{DisableColors: true})
	result.Logger.Level(hlog.DebugLevel)

	OverrideLogger(result.Logger)

	result.Ctx = hlog.Set(context.Background(), result.Logger)
	result.Assert = assert.New(t)
	result.Require = require.New(t)

	return result
}

var oldDefault *hlog.Entry = nil
