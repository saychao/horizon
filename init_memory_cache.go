package horizon

import (
	"github.com/saychao/horizon/cache"
)

func initMemoryCache(app *App) {
	app.cacheProvider = cache.NewProvider()
}

func init() {
	appInit.Add("memory_cache", initMemoryCache)
}
