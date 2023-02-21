package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Config is the configuration for horizon.  It get's populated by the
// app's main function and is provided to NewApp.
type Config struct {
	types.Copuser          `fig:"-"`
	getter                 rawGetter `fig:"-"`
	Hostname               string    `fig:"hostname"`
	DatabaseURL            string    `fig:"database_url,required"`
	StellarCoreDatabaseURL string    `fig:"stellar_core_database_url,required"`
	StellarCoreURL         string    `fig:"stellar_core_url,required"`
	Port                   int       `fig:"port"`

	LogLevel logan.Level `fig:"log_level"`

	SlowQueryBound *time.Duration `fig:"slow_query_bound"`

	APIBackend *url.URL `fig:"api_backend"`

	Ingest bool `fig:"ingest"`
	// StaleThreshold represents the number of ledgers a history database may be
	// out-of-date by before horizon begins to respond with an error to history
	// requests.
	StaleThreshold uint `fig:"stale_threshold"`
	//For developing without signatures
	SkipCheck bool `fig:"sign_checkskip"`
	// enable on dev only
	CORSEnabled bool `fig:"cors_enabled"`
	// DisableAPISubmit tell horizon to not use API for transaction submission
	// for dev purposes only, works well with SkipCheck enabled
	// pending transactions and transaction 2fa will be disabled as well.
	DisableAPISubmit bool `fig:"disable_api_submit"`
	// If set to true - Horizon won't check TFA (via API) during TX submission.
	DisableTXTfa bool `fig:"disable_tx_tfa"`

	TemplateBackend *url.URL `fig:"template_backend"`

	ForceHTTPSLinks bool `fig:"force_https_links"`

	SentryDSN      string `fig:"sentry_dsn"`
	Project        string `fig:"project"`
	SentryLogLevel string `fig:"sentry_log_level"`
	Env            string `fig:"env"`

	MigrateUpOnStart bool `fig:"migrate_up_on_start"`

	CacheSize   int           `fig:"cache_size"`
	CachePeriod time.Duration `fig:"cache_period"`
}

func (c *Config) Init() error {
	err := figure.
		Out(c).
		From(kv.MustGetStringMap(c.getter, "config")).
		With(figure.BaseHooks, URLHook, logLevelHook).
		Please()
	if err != nil {
		return errors.Wrap(err, "failed to figure out config")
	}

	if c.Hostname == "" {
		c.Hostname, err = os.Hostname()
		if err != nil {
			return errors.Wrap(err, "failed to get hostname")
		}
	}
	return nil
}

// rawGetter encapsulates raw config values provider
type rawGetter interface {
	GetStringMap(key string) (map[string]interface{}, error)
}

func NewViperConfig(fn string) Config {
	// init underlying viper
	v := kv.NewViperFile(fn)

	return newViperConfig(v)
}

func newViperConfig(raw rawGetter) Config {
	lvl, err := logan.ParseLevel("debug")
	if err != nil {
		panic("failed to parse default log level")
	}

	config := &Config{
		Port:            8000,
		LogLevel:        lvl,
		ForceHTTPSLinks: true,
		SentryDSN:       "",
		SentryLogLevel:  "warn",
		Project:         "",
		Env:             "",
		SkipCheck:       false,
		Copuser:         copus.NewCopuser(raw),
		CacheSize:       1,
		CachePeriod:     5 * time.Second,
	}

	config.getter = raw
	return *config
}

var (
	URLHook = figure.Hooks{
		"*url.URL": func(value interface{}) (reflect.Value, error) {
			str, err := cast.ToStringE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse string")
			}
			u, err := url.Parse(str)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse url")
			}
			return reflect.ValueOf(u), nil
		},
	}

	logLevelHook = figure.Hooks{
		"map[string]string": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToStringMapStringE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse map[string]string")
			}
			return reflect.ValueOf(result), nil
		},
		"logrus.Level": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				lvl, err := logrus.ParseLevel(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse log level")
				}
				return reflect.ValueOf(lvl), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)
