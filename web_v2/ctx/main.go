package ctx

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"

	"github.com/SafeRE-IT/horizon/config"

	"github.com/SafeRE-IT/horizon/txsub/v2"

	"context"

	"github.com/SafeRE-IT/horizon/corer"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/doorman"
)

type ctxKey int

const (
	keyCoreRepo ctxKey = iota
	keyHistoryRepo
	keyLog
	keyDoorman
	keyCoreInfo
	keyTxSubmitter
	keyConfig
)

// Log - gets entry from context
func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(keyLog).(*logan.Entry)
}

// SetLog - sets log entry into ctx
func SetLog(ctx context.Context, value *logan.Entry) context.Context {
	return context.WithValue(ctx, keyLog, value)
}

// CoreRepo - returns new copy of repo with connection to core DB
func CoreRepo(r *http.Request) *pgdb.DB {
	return getRepo(r, keyCoreRepo)
}

// SetCoreRepo - sets core repo which be used as source for CoreRepo
func SetCoreRepo(repo *pgdb.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyCoreRepo, repo)
	}
}

// HistoryRepo - returns new copy of repo with connection to hisotry DB
func HistoryRepo(r *http.Request) *pgdb.DB {
	return getRepo(r, keyHistoryRepo)
}

// SetHistoryRepo - sets history repo which be used as source for HistoryRepo
func SetHistoryRepo(repo *pgdb.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyHistoryRepo, repo)
	}
}

func getRepo(r *http.Request, key ctxKey) *pgdb.DB {
	repo := r.Context().Value(key).(*pgdb.DB)
	return repo.Clone()
}

// Doorman - perform signature check
func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(keyDoorman).(doorman.Doorman)
	return d.Check(r, constraints...)
}

// SetDoorman - adds doorman to context
func SetDoorman(d doorman.Doorman) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyDoorman, d)
	}
}

// SetCoreInfo - adds core info to context
func SetCoreInfo(coreInfoProvider func() corer.Info) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyCoreInfo, coreInfoProvider())
	}
}

// CoreInfo - returns core info from the context
func CoreInfo(r *http.Request) corer.Info {
	return r.Context().Value(keyCoreInfo).(corer.Info)
}

// Submitter - gets entry from context
func Submitter(r *http.Request) *txsub.System {
	return r.Context().Value(keyTxSubmitter).(*txsub.System)
}

// SetSubmitter - sets submitter entry into ctx
func SetSubmitter(value *txsub.System) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyTxSubmitter, value)
	}
}

// Config - gets entry from context
func Config(r *http.Request) *config.Config {
	return r.Context().Value(keyConfig).(*config.Config)
}

// SetConfig - sets config into ctx
func SetConfig(value *config.Config) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keyConfig, value)
	}
}
