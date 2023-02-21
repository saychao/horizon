package httpx

import (
	"context"

	context2 "github.com/goji/context"
	"github.com/SafeRE-IT/horizon/config"
)

type configgetter interface {
	Conf() config.Config
}

func getConfig(ctx context.Context) config.Config {
	c := context2.ToC(ctx)
	config := c.Env["app"].(configgetter).Conf()
	return config
}
