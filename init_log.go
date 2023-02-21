package horizon

import (
	"time"

	"github.com/evalphobia/logrus_sentry"
	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"fmt"
	"net/http"

	"github.com/saychao/horizon/config"
	"github.com/saychao/horizon/log"
)

const (
	NLinesAroundErrorPoint = 2

	defaultLogLevel = "warn"
)

// initLog initialized the logging subsystem, attaching app.log and
// app.logMetrics.  It also configured the logger's level using Config.LogLevel.
func initLog(app *App) {
	if app.config.SlowQueryBound != nil {
		log.SlowQueryBound = *app.config.SlowQueryBound
	}
	entry := logan.New()
	entry, err := addSentryHook(app.config, entry)
	if err != nil {
		errors.Wrap(err, "Failed to add Sentry hook")
	}
	log.DefaultLogger.Entry = *entry
	log.DefaultLogger.Level(app.config.LogLevel)
	log.DefaultLogger = log.DefaultLogger.WithField("host", app.config.Hostname)

	entry, err = addSentryHook(app.config, entry)
	if err != nil {
		errors.Wrap(err, "Failed to add Sentry hook")
	}
	log.DefaultLogger.Entry = *entry

}

func init() {
	appInit.Add("log", initLog)
}

func addSentryHook(config config.Config, entry *logan.Entry) (*logan.Entry, error) {
	sentry := config.SentryDSN
	lvl := config.SentryLogLevel
	if lvl == "" {
		lvl = defaultLogLevel
	}
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		errors.Wrap(err, "Failed to parse log level")
	}

	if sentry == "" {
		return entry, nil
	}

	hook, err := logrus_sentry.NewSentryHook(sentry, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Sentry hook")
	}

	env := config.Env
	if env == "" {
		env = "unknown"
	}
	hook.SetEnvironment(env)

	proj := config.Project
	if proj == "" {
		proj = "unknown"
	}

	entry = entry.WithField("tags", raven.Tags{
		{
			Key:   "project",
			Value: proj,
		},
	})

	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.Level = level
	hook.StacktraceConfiguration.Context = NLinesAroundErrorPoint
	hook.Timeout = 3 * time.Second
	hook.AddExtraFilter("status_code", func(v interface{}) interface{} {
		i, ok := v.(int)
		if !ok {
			return v
		}

		return fmt.Sprintf("%d - %s", i, http.StatusText(i))
	})

	hook.SetRelease(version)

	wrapperHook := sentryWrapperHook{
		SentryHook: hook,
	}

	entry.AddLogrusHook(&wrapperHook)
	return entry, nil
}

type sentryWrapperHook struct {
	*logrus_sentry.SentryHook
}

func (h *sentryWrapperHook) Fire(entry *logrus.Entry) error {
	err, ok := entry.Data[logan.ErrorKey]
	if ok {
		entry.Data["raw_error"] = err
	}

	entry = h.putServiceToTags(entry)

	return h.SentryHook.Fire(entry)
}

// TODO Make a common helper for field reputting.
func (h *sentryWrapperHook) putServiceToTags(entry *logrus.Entry) *logrus.Entry {
	serviceField, ok := entry.Data["service"]
	if !ok {
		// No 'service' field
		return entry
	}

	serviceName, ok := serviceField.(string)
	if !ok {
		// Service field is not a string
		return entry
	}

	serviceTag := raven.Tag{
		Key:   "service",
		Value: serviceName,
	}

	tagsField, ok := entry.Data["tags"]
	if ok {
		// Try to put service into tags.
		tags, ok := tagsField.(raven.Tags)
		if !ok {
			// Tags field is not a raven.Tags instance. That's quite strange though.
			return entry
		}

		entry.Data["tags"] = append(tags, serviceTag)
	} else {
		// No tags field.
		entry = entry.WithField("tags", raven.Tags{
			serviceTag,
		})
	}

	return entry
}
